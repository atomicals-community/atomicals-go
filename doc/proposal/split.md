### A solution to split atomicals ft

We need a solution to split Atomicals FT, especially non-integer splitting.

### Current Protocol
In the current protocol, each FT is paired with one Satoshi. During transfers, UTXOs that are not fully colored result in the corresponding FT being burned.

For example, Alice initiates a transfer Tx1:
``` 
Vin[0]: 1000 Satoshis (colored with 1000 Atoms)
Vin[1]: 500 Satoshis
The recipients of the transfer include:
Vout[0]: 800 Satoshis (colored with 800 Atoms) to Bob
Vout[1]: 700 Satoshis to Tom
``` 
According to the current Atomicals protocol, UTXOs that cannot be fully colored correspond to illegal FTs. In the above transaction, 200 Atoms are burned.

### Splitting Solution - Supporting Partially Colored UTXOs
In practice, the storage space occupied by the indexer remains the same regardless of whether the UTXO is fully colored or not. The structure of colored UTXOs is defined in atomicals/DB/postsql/atomicalsUTXOFt.go, where the indexer needs to store the ftAmount corresponding to the UTXO, which currently equals the amount of Satoshis. The indexer can record the coloring result of Vout[1] without burned these 200 Atoms.

If the Atomicals protocol recognizes partially colored UTXOs, the result of the above transaction will be:
``` 
Vout[0]: 800 Satoshis (colored with 800 Atoms) to Bob
Vout[1]: 700 Satoshis (colored with 200 Atoms) to Tom
``` 

### Transfer of Partially Colored UTXOs
Tom initiates a new transaction Tx2:
```
Vin[0]: 700 Satoshis (colored with 200 Atoms)
The recipients of the transfer include:
Vout[0]: 500 Satoshis (colored with 200*5/7 Atoms) to Jerry
Vout[1]: 200 Satoshis (colored with 200*2/7 Atoms) to Tom
```
In this transaction, Jerry and Tom will receive a non-integer number of Atoms, with the decimal part depending on the precision specified by the protocol. If we specify support for 10 decimal places, Jerry will receive 142.8571428571 Atoms, and Tom will receive 57.1428571428. Values beyond the 10th decimal place will be destroyed.
