try
{
    clear
    sl cmd/client
    go generate
    go run .
}
finally
{
    sl ../../
}