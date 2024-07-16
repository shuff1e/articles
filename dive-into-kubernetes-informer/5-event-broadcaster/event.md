```markdown
 ~  kubectl get event -o wide abc.17e2a4d750b6fbe0 -o yaml
apiVersion: v1
count: 7
eventTime: null
firstTimestamp: "2024-07-16T08:40:52Z"
involvedObject:
  apiVersion: v1
  kind: Pod
  name: abc
kind: Event
lastTimestamp: "2024-07-16T08:41:52Z"
message: test a b
metadata:
  annotations:
    a: b
  creationTimestamp: "2024-07-16T08:40:52Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:count: {}
      f:firstTimestamp: {}
      f:involvedObject:
        f:apiVersion: {}
        f:kind: {}
        f:name: {}
      f:lastTimestamp: {}
      f:message: {}
      f:metadata:
        f:annotations:
          .: {}
          f:a: {}
      f:reason: {}
      f:source:
        f:component: {}
      f:type: {}
    manager: ___go_build_github_com_wbsnail_articles_archive_dive_into_kubernetes_informer_basic_controller_dive_into_kubernetes_informer_5_e
    operation: Update
    time: "2024-07-16T08:40:52Z"
  name: abc.17e2a4d750b6fbe0
  namespace: default
  resourceVersion: "490472"
  uid: d737ee36-0049-4b8f-97ad-3746fd31be2b
reason: Evicted
reportingComponent: ""
reportingInstance: ""
source:
  component: test-leader-election
type: Warning
```