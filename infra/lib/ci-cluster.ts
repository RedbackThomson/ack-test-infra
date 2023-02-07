import { Construct } from "constructs";
import { aws_eks as eks, aws_ec2 as ec2, aws_iam as iam } from "aws-cdk-lib";
import * as cdk8s from "cdk8s";
import { policies as ALBPolicies } from "./policies/aws-load-balancer-controller-policy";
import {
  ProwGitHubSecretsChart,
  ProwGitHubSecretsChartProps,
} from "./charts/prow-secrets";
import {
  EXTERNAL_DNS_NAMESPACE,
  FLUX_NAMESPACE,
  PROW_JOB_NAMESPACE,
  PROW_NAMESPACE,
} from "./test-ci-stack";

export type CIClusterCompileTimeProps = ProwGitHubSecretsChartProps;

export type CIClusterRuntimeProps = {};

export type CIClusterProps = CIClusterCompileTimeProps & CIClusterRuntimeProps;

export class CICluster extends Construct {
  readonly testCluster: eks.Cluster;
  readonly testNodegroup: eks.Nodegroup;

  readonly namespaceManifests: eks.KubernetesManifest[];

  constructor(scope: Construct, id: string, props: CIClusterProps) {
    super(scope, id);

    this.testCluster = new eks.Cluster(scope, "TestInfraCluster", {
      version: eks.KubernetesVersion.V1_24,
      defaultCapacity: 0,
    });
    this.testNodegroup = this.testCluster.addNodegroupCapacity(
      "TestInfraNodegroup",
      {
        instanceTypes: [
          ec2.InstanceType.of(ec2.InstanceClass.M5, ec2.InstanceSize.XLARGE8),
        ],
        minSize: 2,
        diskSize: 150,
      }
    );

    this.namespaceManifests = [
      EXTERNAL_DNS_NAMESPACE,
      PROW_JOB_NAMESPACE,
      PROW_NAMESPACE,
    ].map(this.createNamespace);

    this.installProwRequirements(props);
    this.installFlux();
    this.installExternalDNS();
    this.installAWSLoadBalancer();
  }

  createNamespace = (name: string) => {
    return new eks.KubernetesManifest(
      this.testCluster.stack,
      `${name}-namespace-struct`,
      {
        cluster: this.testCluster,
        manifest: [
          {
            apiVersion: "v1",
            kind: "Namespace",
            metadata: {
              name: name,
            },
          },
        ],
      }
    );
  };

  installFlux = () => {
    const fluxChart = this.testCluster.addHelmChart("flux2", {
      chart: "flux2",
      repository: "https://fluxcd-community.github.io/helm-charts",
      namespace: FLUX_NAMESPACE,
      createNamespace: true,
      version: "0.19.2",
      values: {},
    });

    const fluxBootstrap = this.testCluster.addManifest(
      "FluxBootstrap",
      ...[
        {
          apiVersion: "source.toolkit.fluxcd.io/v1beta2",
          kind: "GitRepository",
          metadata: {
            name: "test-infra",
            namespace: "flux-system",
          },
          spec: {
            interval: "30s",
            ref: {
              branch: "main",
            },
            url: "https://github.com/aws-controllers-k8s/test-infra",
          },
        },
        {
          apiVersion: "kustomize.toolkit.fluxcd.io/v1beta2",
          kind: "Kustomization",
          metadata: {
            name: "all-apps",
            namespace: "flux-system",
          },
          spec: {
            interval: "5m",
            sourceRef: {
              kind: "GitRepository",
              name: "test-infra",
            },
            path: "./flux",
            prune: true,
            targetNamespace: "flux-system",
            validation: "client",
          },
        },
      ]
    );
    fluxBootstrap.node.addDependency(fluxChart);
  };

  installProwRequirements = (secretsProps: ProwGitHubSecretsChartProps) => {
    const prowSecretsApp = new cdk8s.App();
    const prowSecretsChart = this.testCluster.addCdk8sChart(
      "prow-secrets",
      new ProwGitHubSecretsChart(prowSecretsApp, "ProwSecrets", secretsProps)
    );

    // Ensure namespaces are created before secrets
    prowSecretsChart.node.addDependency(...this.namespaceManifests);
    prowSecretsApp.charts.forEach((chart) =>
      chart.addDependency(...this.namespaceManifests)
    );
  };

  installExternalDNS = () => {
    const externalDNSServiceAccount = this.testCluster.addServiceAccount(
      "external-dns-service-account",
      {
        namespace: EXTERNAL_DNS_NAMESPACE,
      }
    );
    externalDNSServiceAccount.node.addDependency(...this.namespaceManifests);

    externalDNSServiceAccount.addToPrincipalPolicy(
      new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: ["route53:ChangeResourceRecordSets"],
        resources: ["arn:aws:route53:::hostedzone/*"],
      })
    );
    externalDNSServiceAccount.addToPrincipalPolicy(
      new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: ["route53:ListHostedZones", "route53:ListResourceRecordSets"],
        resources: ["*"],
      })
    );

    const helmChart = this.testCluster.addHelmChart("external-dns", {
      chart: "external-dns",
      repository: "https://charts.bitnami.com/bitnami",
      namespace: EXTERNAL_DNS_NAMESPACE,
      version: "6.12.0",
      values: {
        namespace: PROW_NAMESPACE, // Limit only to DNS in Prow
        sources: ["ingress"],
        policy: "upsert-only",
        serviceAccount: {
          create: false,
          name: externalDNSServiceAccount.serviceAccountName,
        },
        aws: {
          zoneType: "public",
        },
      },
    });
    helmChart.node.addDependency(...this.namespaceManifests);
  };

  installAWSLoadBalancer = () => {
    const serviceAccount = this.testCluster.addServiceAccount(
      "alb-service-account",
      {
        namespace: "kube-system",
      }
    );
    ALBPolicies.map((policy) => serviceAccount.addToPrincipalPolicy(policy));

    this.testCluster.addHelmChart("aws-load-balancer-controller", {
      chart: "aws-load-balancer-controller",
      repository: "https://aws.github.io/eks-charts",
      namespace: "kube-system",
      version: "1.4.7",
      values: {
        clusterName: this.testCluster.clusterName,
        serviceAccount: {
          create: false,
          name: serviceAccount.serviceAccountName,
        },
      },
    });
  };
}
