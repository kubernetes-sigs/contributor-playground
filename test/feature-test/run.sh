#!/bin/bash

startTime=`date +%Y%m%d-%H:%M`
startTime_s=`date +%s`

echo "Test nginx.yaml"
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
echo "Test nginx.yaml success"
echo ""

echo "Test nginx-BLB-allocate-vip.yaml"
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
echo "Test nginx-BLB-allocate-vip.yaml success"
echo ""

echo "Test nginx-BLB-support-internal-vpc.yaml"
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
echo "Test nginx-BLB-support-internal-vpc.yaml success"
echo ""

echo "Test nginx-EIP-Postpaid-ByBandwidth-50M.yaml"
kubectl create -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByBandwidth-50M.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-eip-postpaid-by-bandwidth-50m"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByBandwidth-50M.yaml
echo "Test nginx-EIP-Postpaid-ByBandwidth-50M.yaml success"
echo ""

echo "Test nginx-EIP-Postpaid-ByBandwidth-200M.yaml"
kubectl create -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByBandwidth-200M.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-eip-postpaid-by-bandwidth-200m"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByBandwidth-200M.yaml
echo "Test nginx-EIP-Postpaid-ByBandwidth-200M.yaml success"
echo ""

echo "Test nginx-EIP-Postpaid-ByTraffic-50M.yaml"
kubectl create -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByTraffic-50M.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-eip-postpaid-by-traffic-50m"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByTraffic-50M.yaml
echo "Test nginx-EIP-Postpaid-ByTraffic-50M.yaml success"
echo ""

echo "Test nginx-EIP-Postpaid-ByTraffic-1000M.yaml"
kubectl create -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByTraffic-1000M.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-eip-postpaid-by-traffic-1000m"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
kubectl delete -f ./docs/controllers/service/examples/nginx-EIP-Postpaid-ByTraffic-1000M.yaml
echo "Test nginx-EIP-Postpaid-ByTraffic-1000M.yaml success"
echo ""

echo "Test nginx-BLB-assignedID.yaml"
kubectl create -f ./docs/controllers/service/examples/nginx-BLB-assignedID.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-blb-assigned-id"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
echo "echo create successfully"
sleep 10
kubectl delete -f ./docs/controllers/service/examples/nginx-BLB-assignedID.yaml
echo "Test nginx-BLB-assignedID.yaml successfully"
echo ""

echo "Test nginx-BLB-assignedID-EIP.yaml"
kubectl create -f ./docs/controllers/service/examples/nginx-BLB-assignedID-EIP.yaml
status="Check: "
while true
do
    if [ `kubectl get service |grep "nginx-service-blb-assignedid-eip"|grep -v "pending"|wc -l` -eq 1 ]; then
        break
    fi
    sleep 10
    status=$status"."
    echo $status
done
echo "echo create successfully"
sleep 10
kubectl delete -f ./docs/controllers/service/examples/nginx-BLB-assignedID-EIP.yaml
echo "Test nginx-BLB-assignedID-EIP.yaml successfully"
echo ""


endTime=`date +%Y%m%d-%H:%M`
endTime_s=`date +%s`
sumTime=$[ $endTime_s - $startTime_s ]
echo "Test Finish: feature-test" "$startTime ---> $endTime" "Total: $sumTime s"
echo ""