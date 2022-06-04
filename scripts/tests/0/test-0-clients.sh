for i in `seq 1 3`
do
  make restart
  start=$(date +%s)
  for j in `seq 1 40`
  do
    current=$(date +%s)
    echo "$(($current-$start)) $(docker stats mac0352-ep2_server_1 --format "{{.CPUPerc}}" --no-stream)" >> "$1-$i"
    echo "$(($current-$start)) $(docker stats mac0352-ep2_server_1 --format "{{.NetIO}}" --no-stream)" >> "$2-$i"
  done
done

make shutdown