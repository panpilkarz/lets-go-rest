## Test

```
docker-compose up
```

## Expected output

```
$ docker-compose up | grep accountsclienttest

accountsclienttest_1  | === RUN   TestAllOperations
accountsclienttest_1  | --- PASS: TestAllOperations (0.03s)
accountsclienttest_1  | === RUN   TestInvalidEndpoint
accountsclienttest_1  | --- PASS: TestInvalidEndpoint (0.00s)
accountsclienttest_1  | === RUN   TestCreate
accountsclienttest_1  | --- PASS: TestCreate (0.01s)
accountsclienttest_1  | === RUN   TestFetch
accountsclienttest_1  | --- PASS: TestFetch (0.01s)
accountsclienttest_1  | === RUN   TestDelete
accountsclienttest_1  | --- PASS: TestDelete (0.01s)
accountsclienttest_1  | === RUN   TestList
accountsclienttest_1  | --- PASS: TestList (0.08s)
accountsclienttest_1  | === RUN   TestListPage
accountsclienttest_1  | --- PASS: TestListPage (0.08s)
accountsclienttest_1  | PASS
accountsclienttest_1  | ok      accountsclient  0.223s
```
