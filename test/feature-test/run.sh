#!/bin/bash

startTime=`date +%Y%m%d-%H:%M`
startTime_s=`date +%s`

echo "Test Nginx"
kubectl create -f ./docs/controllers/service/examples/nginx.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx.yaml
echo "Test Nginx success"
echo ""

echo "Test BLB allocate vip"
kubectl create -f ./docs/controllers/service/examples/nginx-BLB-allocate-vip.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-blb-allocate-vip"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx-BLB-allocate-vip.yaml
echo "Test BLB allocate vip success"
echo ""

echo "Test BLB support internal vpc"
kubectl create -f ./docs/controllers/service/examples/nginx-BLB-support-internal-vpc.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-blb-internal-vpc"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx-BLB-support-internal-vpc.yaml
echo "Test BLB support internal vpc success"
echo ""

endTime=`date +%Y%m%d-%H:%M`
endTime_s=`date +%s`
sumTime=$[ $endTime_s - $startTime_s ]
echo "Test Finish: feature-test" "$startTime ---> $endTime" "Total: $sumTime s"