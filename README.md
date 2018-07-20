Personal reminder for how to start geth
nohup geth --datadir ./myDataDir --verbosity=1 --mine --minerthreads=1 --rpc --rpcaddr 0.0.0.0 --rpcapi "admin,eth,miner,net,personal,web3" &

Getting go-ethereum dependencies are somewhat of a paid to build.  You need a C compiler in your path with CGO_ENABLED environment variable set to 1.
If you have trouble building, create an issue and I will attempt to help.