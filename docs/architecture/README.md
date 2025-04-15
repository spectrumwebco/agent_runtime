# ML Infrastructure Architecture

This document describes the architecture of the ML infrastructure for fine-tuning and evaluating Llama 4 models.

## System Architecture

The ML infrastructure follows a modular architecture with components deployed on Kubernetes. The system is designed to be scalable, resilient, and easy to maintain.

```mermaid
graph TD
    subgraph Data Sources
        GH[GitHub Issues]
        GE[Gitee Issues]
    end

    subgraph Data Pipeline
        SC[Scraper Connector]
        DP[Data Preprocessor]
        DV[Data Validator]
        DT[Dataset Versioner]
    end

    subgraph Feature Store
        FS[Feast]
    end

    subgraph Training Infrastructure
        KF[KubeFlow]
        ML[MLFlow]
        JL[JupyterLab]
        H2O[h2o.ai]
    end

    subgraph Model Serving
        KS[KServe]
        SD[Seldon Core]
    end

    subgraph Storage
        MO[MinIO]
    end

    GH --> SC
    GE --> SC
    SC --> DP
    DP --> DV
    DV --> DT
    DT --> FS
    DT --> MO
    FS --> KF
    MO --> KF
    MO --> ML
    KF --> ML
    ML --> KS
    ML --> SD
    JL --> KF
    H2O --> KF
```

## Component Architecture

### Data Pipeline

The data pipeline collects, processes, and prepares training data from GitHub and Gitee repositories.

```mermaid
graph TD
    subgraph Data Collection
        GH[GitHub API]
        GE[Gitee API]
        SC[Scraper Connector]
    end

    subgraph Data Processing
        DP[Data Preprocessor]
        DV[Data Validator]
        DT[Dataset Versioner]
    end

    subgraph Storage
        MO[MinIO]
    end

    GH --> SC
    GE --> SC
    SC --> DP
    DP --> DV
    DV --> DT
    DT --> MO
```

### Training Infrastructure

The training infrastructure orchestrates the fine-tuning of Llama 4 models.

```mermaid
graph TD
    subgraph KubeFlow
        KFP[Pipelines]
        KFT[Training Operator]
        KFK[Katib]
        KFN[Notebooks]
    end

    subgraph MLFlow
        MLE[Experiment Tracking]
        MLM[Model Registry]
    end

    subgraph Storage
        MO[MinIO]
    end

    subgraph Compute
        GPU[GPU Nodes]
    end

    KFP --> KFT
    KFP --> KFK
    KFT --> GPU
    KFK --> GPU
    KFN --> KFP
    KFT --> MLE
    MLE --> MLM
    MO --> KFP
    MLE --> MO
    MLM --> MO
```

### Model Serving

The model serving infrastructure deploys and serves trained models.

```mermaid
graph TD
    subgraph KServe
        KSI[InferenceService]
        KSC[Canary Deployment]
        KSS[Scaling]
    end

    subgraph Seldon Core
        SDP[Deployment]
        SIP[Inference Pipeline]
    end

    subgraph Storage
        MO[MinIO]
    end

    subgraph MLFlow
        MLM[Model Registry]
    end

    MLM --> KSI
    MLM --> SDP
    MO --> KSI
    MO --> SDP
    KSI --> KSC
    KSI --> KSS
    SDP --> SIP
```

## Deployment Architecture

The ML infrastructure is deployed on Kubernetes using a combination of Kubernetes manifests and Terraform configurations.

```mermaid
graph TD
    subgraph Terraform
        TFM[Modules]
        TFP[Providers]
        TFV[Variables]
    end

    subgraph Kubernetes
        KNS[Namespaces]
        KDP[Deployments]
        KSV[Services]
        KIN[Ingress]
        KPV[PersistentVolumes]
        KCM[ConfigMaps]
        KSC[Secrets]
    end

    TFM --> KNS
    TFM --> KDP
    TFM --> KSV
    TFM --> KIN
    TFM --> KPV
    TFM --> KCM
    TFM --> KSC
    TFP --> TFM
    TFV --> TFM
```

## Data Flow

The following diagram illustrates the data flow through the ML infrastructure.

```mermaid
graph LR
    subgraph Data Sources
        GH[GitHub Issues]
        GE[Gitee Issues]
    end

    subgraph Data Pipeline
        SC[Scraper Connector]
        DP[Data Preprocessor]
        DV[Data Validator]
        DT[Dataset Versioner]
    end

    subgraph Storage
        MO[MinIO]
    end

    subgraph Training
        KF[KubeFlow]
        ML[MLFlow]
    end

    subgraph Serving
        KS[KServe]
        SD[Seldon Core]
    end

    subgraph Clients
        AP[API Clients]
        WB[Web UI]
    end

    GH --> SC
    GE --> SC
    SC --> DP
    DP --> DV
    DV --> DT
    DT --> MO
    MO --> KF
    KF --> ML
    ML --> MO
    MO --> KS
    MO --> SD
    KS --> AP
    SD --> AP
    KS --> WB
    SD --> WB
```

## Security Architecture

The ML infrastructure implements security at multiple levels.

```mermaid
graph TD
    subgraph Authentication
        AU[User Authentication]
        AK[API Keys]
        AT[Service Tokens]
    end

    subgraph Authorization
        RB[RBAC]
        NS[Namespaces]
        NP[Network Policies]
    end

    subgraph Encryption
        ET[TLS]
        ES[Secret Management]
    end

    AU --> RB
    AK --> RB
    AT --> RB
    RB --> NS
    NS --> NP
    ET --> AU
    ES --> AU
    ES --> AK
    ES --> AT
```
