- [2020 Q1 (Jan-Mar)](#sec-1)
  - [Increase Stable Test Coverage Velocity 100% over 2019 (Score: 0.9)](#sec-1-1)
    - [KR1=0.9 (26/+27) new conformant stable endpoints](#sec-1-1-1)
    - [KR2=0.9 +6% Coverage Increase](#sec-1-1-2)
  - [Complete cncf/apisnoop prow.k8s.io + EKS migration (Score: 0.5)](#sec-1-2)
    - [KR1=0.5 All cncf/apisnoop artifacts created by prow.k8s.io](#sec-1-2-1)
    - [KR2=0.0 All cncf/apisnoop github workflow managed by prow.k8s.io](#sec-1-2-2)
    - [KR3=1.0 All cncf/apisnoop non-prow infra moved to EKS](#sec-1-2-3)
  - [Mentor/Teach test-writing workflow at Contributer Summit / KubeConEU (Score: 1.0)](#sec-1-3)
    - [KR1 Caleb and Hippie Mentoring at Contributor Summit](#sec-1-3-1)
    - [KR2 Zach and Stephen teaching test writing](#sec-1-3-2)
- [2020 Q2 (Apr-Jun)](#sec-2)
  - [Increase Stable Test Coverage Velocity 50% over Q1](#sec-2-1)
    - [KR1 (0/+40) new conformant stable endpoints](#sec-2-1-1)
    - [KR2 +9% Coverage Increase](#sec-2-1-2)
    - [KR3 (stretch) 50% stable endpoints hit by conformance tests](#sec-2-1-3)
  - [Prepare to Gate k/k PRs touching test/e2e or API](#sec-2-2)
    - [KR1 comment w/ list of increase/decrease of stable endpoints](#sec-2-2-1)
    - [KR2 gate w/ comment](#sec-2-2-2)
  - [Prepare to Gate cncf/k8s-conformance PRs touching v\*.\*/](#sec-2-3)
    - [KR1 comment w/ list of unrun conformance tests](#sec-2-3-1)
    - [KR2 gate w/ comment](#sec-2-3-2)


# 2020 Q1 (Jan-Mar)<a id="sec-1"></a>

## Increase Stable Test Coverage Velocity 100% over 2019 (Score: 0.9)<a id="sec-1-1"></a>

### KR1=0.9 (26/+27) new conformant stable endpoints<a id="sec-1-1-1"></a>

1.  SCORE: 0.9

    Done = 5 Needs Approval = 3 Needs Review = 9 In Progress (no flakes) = 9

2.  Done = 5

    1.  DONE +3 Promote: Secret patching test #87262
    
    2.  DONE +1 Promote: find Kubernetes Service in default Namespace #87260
    
    3.  DONE +1 Promote: Namespace patch test #87256

3.  Needs Approval +3

    1.  PROMOTION +3 Promote: pod PreemptionExecutionPath verification
    
        -   ? #issue
        -   ? #test
        -   Promotion: <https://github.com/kubernetes/kubernetes/pull/83378>
        
        Clayton says: "I got it" Has a failing test&#x2026; /retest

4.  Needs Review +9

    1.  TEST +3 Promote: PodTemplate Lifecycle test #88036 (removing flakes)
    
        -   Issue: <https://github.com/kubernetes/kubernetes/issues/86141> #issue Needs reopening and checkboxes for current state..
        -   Promotion: <https://github.com/kubernetes/kubernetes/pull/88036#ref-pullrequest-571656281>
        -   Flakes: <https://github.com/kubernetes/kubernetes/pull/88588#issuecomment-606957802>
        -   Addressing Flakes: [https://github.com/kubernetes/kubernetes/pull/89746](https://github.com/kubernetes/kubernetes/pull/89746)
    
    2.  COMMENTS +2 Promote: ConfigMap Lifecycle test #88034 (comments addressed)
    
        -   Promotion: <https://github.com/kubernetes/kubernetes/pull/88034#discussion_r398728147>
        -   Addressing Comments: <https://github.com/kubernetes/kubernetes/pull/88034#issuecomment-607430447> (addresed)
        -   PR to handle timeouts: <https://github.com/kubernetes/kubernetes/pull/89707>
    
    3.  COMMENTS +4 Pod and PodStatus
    
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/88545>
        -   test: <https://github.com/kubernetes/kubernetes/pull/89453> Addressed the [comment](https://github.com/kubernetes/kubernetes/pull/89453#discussion_r400346746): "Not sure this will work, you will be racing with the kubelet, I think. That is, kubelet may mark it ready again."

5.  In Progress +18

    1.  TEST +4 Promote: Endpoints
    
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/87762>
        -   test: <https://github.com/kubernetes/kubernetes/pull/88778>
        -   promotion: <https://github.com/kubernetes/kubernetes/pull/89752> [TestGrid reference](https://testgrid.k8s.io/sig-release-master-blocking#gce-cos-master-default&include-filter-by-regex=should%20test%20the%20lifecycle%20of%20an%20Endpoint) still looks green!
    
    2.  TEST +5 Promote: Event Lifecycle test #86858
    
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/86288>
        -   test: <https://github.com/kubernetes/kubernetes/pull/86858>
        -   promotion: <https://github.com/kubernetes/kubernetes/pull/89753> [TestGrid reference](https://testgrid.k8s.io/sig-release-master-blocking#gce-cos-master-default&include-filter-by-regex=should%20ensure%20that%20an%20event%20can%20be%20fetched%2C%20patched%2C%20deleted%2C%20and%20listed)
    
    3.  FLAKING +7 ReplicationController lifecycle
    
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/88302> Needs reopening and checkboxes for current state&#x2026;
        -   test: <https://github.com/kubernetes/kubernetes/pull/88588>
        -   [address flaking comment](https://github.com/kubernetes/kubernetes/issues/89740) : [https://github.com/kubernetes/kubernetes/pull/89746](https://github.com/kubernetes/kubernetes/pull/89746)
        
        relies on it's own update response data

6.  Sorted Backlog +5

    1.  BACKLOG +2 ServiceStatus lifecycle
    
        -   org-file: <https://github.com/cncf/apisnoop/pull/298>
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/89135> Currently, this test is having issues writing to the ServiceStatus endpoints (via patch and update). The data is patched without errors, but the data when fetched is no different to before the patching.
    
    2.  BACKLOG +3 ServiceAccount lifecycle
    
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/89071> @johnbelamaric You don't need to check the status of the secret as part of the test. In other places we check that the resource in question happens, we don't have to follow.

7.  Triage +12

    1.  TRIAGE +5 Apps DaemonSet lifecycle
    
        -   org-file: <https://github.com/cncf/apisnoop/pull/305>
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/89637>
    
    2.  TRIAGE +5 Apps Deployment lifecycle
    
        -   org-file:
        -   mock-test: <https://github.com/kubernetes/kubernetes/issues/89340>
    
    3.  TRIAGE +2 NodeStatus     :deprioritized:
    
        Needs these comments addressed, and we voted to de-priorize <https://github.com/kubernetes/kubernetes/issues/88358#issuecomment-591062171>

### KR2=0.9 +6% Coverage Increase<a id="sec-1-1-2"></a>

1.  SCORE: 0.9

    Based on the same KR above&#x2026; it's not merged, but these are non-flakey tests that are ready to merge.

## Complete cncf/apisnoop prow.k8s.io + EKS migration (Score: 0.5)<a id="sec-1-2"></a>

### KR1=0.5 All cncf/apisnoop artifacts created by prow.k8s.io<a id="sec-1-2-1"></a>

### KR2=0.0 All cncf/apisnoop github workflow managed by prow.k8s.io<a id="sec-1-2-2"></a>

### KR3=1.0 All cncf/apisnoop non-prow infra moved to EKS<a id="sec-1-2-3"></a>

## Mentor/Teach test-writing workflow at Contributer Summit / KubeConEU (Score: 1.0)<a id="sec-1-3"></a>

### KR1 Caleb and Hippie Mentoring at Contributor Summit<a id="sec-1-3-1"></a>

I am pairing weekly with Guin and Mallian to ensure the workflow is accessible.

### KR2 Zach and Stephen teaching test writing<a id="sec-1-3-2"></a>

They are teaching Riaan, all remote, using our org-flow.

# 2020 Q2 (Apr-Jun)<a id="sec-2"></a>

## Increase Stable Test Coverage Velocity 50% over Q1<a id="sec-2-1"></a>

### KR1 (0/+40) new conformant stable endpoints<a id="sec-2-1-1"></a>

### KR2 +9% Coverage Increase<a id="sec-2-1-2"></a>

### KR3 (stretch) 50% stable endpoints hit by conformance tests<a id="sec-2-1-3"></a>

## Prepare to Gate k/k PRs touching test/e2e or API<a id="sec-2-2"></a>

### KR1 comment w/ list of increase/decrease of stable endpoints<a id="sec-2-2-1"></a>

### KR2 gate w/ comment<a id="sec-2-2-2"></a>

## Prepare to Gate cncf/k8s-conformance PRs touching v\*.\*/<a id="sec-2-3"></a>

### KR1 comment w/ list of unrun conformance tests<a id="sec-2-3-1"></a>

### KR2 gate w/ comment<a id="sec-2-3-2"></a>
