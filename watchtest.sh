while true
do
    echo '--------------------------------------------------------------------------------'
    inotifywait -qq -r -e create,close_write,modify,move,delete ./ && go test ./...

done
