


# Sharder API:
  

## Informations

### Version

0.1.0

## Content negotiation

### URI Schemes
  * http
  * https

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  miner

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /v1/block/get/fee_stats | [get block fee stats](#get-block-fee-stats) | Get block fee stats. |
| GET | /v1/diagnostics/get/info | [get diagnostics info](#get-diagnostics-info) | Get latest block and round metrics cached in the miner. |
| GET | /v1/block/get/latest_finalized | [get latest finalized block](#get-latest-finalized-block) | Get latest finalized block. |
| GET | /v1/block/get/recent_finalized | [get recent finalized block](#get-recent-finalized-block) | Get recent finalized blocks. |
| GET | /v1/estimate_txn_fee | [get txn fees](#get-txn-fees) |  |
| GET | /v1/fees_table | [get txn fees table](#get-txn-fees-table) |  |
  


###  sharder

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /v1/block/get | [get block](#get-block) | Get Block. |
| GET | /v1/client/get/balance | [get client balance](#get-client-balance) | Get client balance. |
| GET | /v1/current-round | [get current round](#get-current-round) | Get round. |
| GET | /v1/healthcheck | [get health check](#get-health-check) | Health Check. |
| GET | /v1/block/magic/get | [get magic block](#get-magic-block) | Get Magic Block. |
| GET | /v1/sharder/get/stats | [get sharder stats](#get-sharder-stats) | Get Sharder Stats. |
| GET | /v1/transaction/get/confirmation | [get transaction confirmationz](#get-transaction-confirmationz) | Get transaction confirmation. |
  


## Paths

### <span id="get-block"></span> Get Block. (*GetBlock*)

```
GET /v1/block/get
```

Retrieve needed parts of block information, given either its round or its hash. At least one of them needs to be provided, if both are provided, however, the round will overwrite the hash.
If "content" == "full", the response has the full Block in `block` field.
If "content" == "header", the response has the BlockSummary in `header` field.
If "content" == "merkle_tree", the response has the Merkle Tree of the transactions in the block in `merkle_tree` field.

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| block | `query` | string | `string` |  |  |  | Hash of the block to retrieve. |
| content | `query` | string | `string` |  |  |  | A comma-separated list of parts of the block to retrieve. Possible values are "full" to retrieve the full block, "header" to retrieve summary, "merkle_tree" to retrieve Merkle Tree of the transactions inside the block. Default is "header". |
| round | `query` | string | `string` |  |  |  | Round of the block to retrieve. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-block-200) | OK | BlockResponse |  | [schema](#get-block-200-schema) |
| [400](#get-block-400) | Bad Request |  |  | [schema](#get-block-400-schema) |

#### Responses


##### <span id="get-block-200"></span> 200 - BlockResponse
Status: OK

###### <span id="get-block-200-schema"></span> Schema
   
  

[BlockResponse](#block-response)

##### <span id="get-block-400"></span> 400
Status: Bad Request

###### <span id="get-block-400-schema"></span> Schema

### <span id="get-block-fee-stats"></span> Get block fee stats. (*GetBlockFeeStats*)

```
GET /v1/block/get/fee_stats
```

Returns the fee statistics for the transactions of the LFB (latest finalized block). No parameters needed.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-block-fee-stats-200) | OK | BlockFeeStatsResponse |  | [schema](#get-block-fee-stats-200-schema) |

#### Responses


##### <span id="get-block-fee-stats-200"></span> 200 - BlockFeeStatsResponse
Status: OK

###### <span id="get-block-fee-stats-200-schema"></span> Schema
   
  

[BlockFeeStatsResponse](#block-fee-stats-response)

### <span id="get-client-balance"></span> Get client balance. (*GetClientBalance*)

```
GET /v1/client/get/balance
```

Retrieves the balance of a client.

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| client_id | `query` | string | `string` |  | ✓ |  | Client ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-client-balance-200) | OK | State |  | [schema](#get-client-balance-200-schema) |
| [400](#get-client-balance-400) | Bad Request |  |  | [schema](#get-client-balance-400-schema) |

#### Responses


##### <span id="get-client-balance-200"></span> 200 - State
Status: OK

###### <span id="get-client-balance-200-schema"></span> Schema
   
  

[State](#state)

##### <span id="get-client-balance-400"></span> 400
Status: Bad Request

###### <span id="get-client-balance-400-schema"></span> Schema

### <span id="get-current-round"></span> Get round. (*GetCurrentRound*)

```
GET /v1/current-round
```

Retrieves the current round number as int64.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-current-round-200) | OK |  |  | [schema](#get-current-round-200-schema) |
| [400](#get-current-round-400) | Bad Request |  |  | [schema](#get-current-round-400-schema) |

#### Responses


##### <span id="get-current-round-200"></span> 200
Status: OK

###### <span id="get-current-round-200-schema"></span> Schema

##### <span id="get-current-round-400"></span> 400
Status: Bad Request

###### <span id="get-current-round-400-schema"></span> Schema

### <span id="get-diagnostics-info"></span> Get latest block and round metrics cached in the miner. (*GetDiagnosticsInfo*)

```
GET /v1/diagnostics/get/info
```

Returns the latest block/round information known to the node. No parameters needed.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-diagnostics-info-200) | OK | InfoResponse |  | [schema](#get-diagnostics-info-200-schema) |

#### Responses


##### <span id="get-diagnostics-info-200"></span> 200 - InfoResponse
Status: OK

###### <span id="get-diagnostics-info-200-schema"></span> Schema
   
  

[InfoResponse](#info-response)

### <span id="get-health-check"></span> Health Check. (*GetHealthCheck*)

```
GET /v1/healthcheck
```

Retrieve the health check information of the sharder.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-health-check-200) | OK | HealthCheckResponse |  | [schema](#get-health-check-200-schema) |
| [400](#get-health-check-400) | Bad Request |  |  | [schema](#get-health-check-400-schema) |

#### Responses


##### <span id="get-health-check-200"></span> 200 - HealthCheckResponse
Status: OK

###### <span id="get-health-check-200-schema"></span> Schema
   
  

[HealthCheckResponse](#health-check-response)

##### <span id="get-health-check-400"></span> 400
Status: Bad Request

###### <span id="get-health-check-400-schema"></span> Schema

### <span id="get-latest-finalized-block"></span> Get latest finalized block. (*GetLatestFinalizedBlock*)

```
GET /v1/block/get/latest_finalized
```

Retrieves the latest finalized block. No parameters needed.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-latest-finalized-block-200) | OK | BlockSummary |  | [schema](#get-latest-finalized-block-200-schema) |

#### Responses


##### <span id="get-latest-finalized-block-200"></span> 200 - BlockSummary
Status: OK

###### <span id="get-latest-finalized-block-200-schema"></span> Schema
   
  

[BlockSummary](#block-summary)

### <span id="get-magic-block"></span> Get Magic Block. (*GetMagicBlock*)

```
GET /v1/block/magic/get
```

Retrieve the magic block given its number.

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| magic_block_number | `query` | string | `string` |  |  |  | Number of the magic block to retrieve. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-magic-block-200) | OK | Block |  | [schema](#get-magic-block-200-schema) |
| [400](#get-magic-block-400) | Bad Request |  |  | [schema](#get-magic-block-400-schema) |

#### Responses


##### <span id="get-magic-block-200"></span> 200 - Block
Status: OK

###### <span id="get-magic-block-200-schema"></span> Schema
   
  

[Block](#block)

##### <span id="get-magic-block-400"></span> 400
Status: Bad Request

###### <span id="get-magic-block-400-schema"></span> Schema

### <span id="get-recent-finalized-block"></span> Get recent finalized blocks. (*GetRecentFinalizedBlock*)

```
GET /v1/block/get/recent_finalized
```

Returns a list of the 10 most recent finalized blocks. No parameters needed.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-recent-finalized-block-200) | OK | BlockSummary |  | [schema](#get-recent-finalized-block-200-schema) |
| [400](#get-recent-finalized-block-400) | Bad Request |  |  | [schema](#get-recent-finalized-block-400-schema) |

#### Responses


##### <span id="get-recent-finalized-block-200"></span> 200 - BlockSummary
Status: OK

###### <span id="get-recent-finalized-block-200-schema"></span> Schema
   
  

[][BlockSummary](#block-summary)

##### <span id="get-recent-finalized-block-400"></span> 400
Status: Bad Request

###### <span id="get-recent-finalized-block-400-schema"></span> Schema

### <span id="get-sharder-stats"></span> Get Sharder Stats. (*GetSharderStats*)

```
GET /v1/sharder/get/stats
```

Retrieve the sharder stats.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-sharder-stats-200) | OK | ExplorerStats |  | [schema](#get-sharder-stats-200-schema) |
| [404](#get-sharder-stats-404) | Not Found |  |  | [schema](#get-sharder-stats-404-schema) |

#### Responses


##### <span id="get-sharder-stats-200"></span> 200 - ExplorerStats
Status: OK

###### <span id="get-sharder-stats-200-schema"></span> Schema
   
  

[ExplorerStats](#explorer-stats)

##### <span id="get-sharder-stats-404"></span> 404
Status: Not Found

###### <span id="get-sharder-stats-404-schema"></span> Schema

### <span id="get-transaction-confirmationz"></span> Get transaction confirmation. (*GetTransactionConfirmationz*)

```
GET /v1/transaction/get/confirmation
```

Get the confirmation of the transaction from the sharders.
If content == confirmation, only the confirmation is returned. Otherwise, the confirmation and the latest finalized block are returned.

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| content | `query` | string | `string` |  |  | `"confirmation"` | confirmation or error |
| hash | `query` | string | `string` |  | ✓ |  | Transaction hash |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-transaction-confirmationz-200) | OK | ConfirmationResponse |  | [schema](#get-transaction-confirmationz-200-schema) |
| [400](#get-transaction-confirmationz-400) | Bad Request |  |  | [schema](#get-transaction-confirmationz-400-schema) |

#### Responses


##### <span id="get-transaction-confirmationz-200"></span> 200 - ConfirmationResponse
Status: OK

###### <span id="get-transaction-confirmationz-200-schema"></span> Schema
   
  

[ConfirmationResponse](#confirmation-response)

##### <span id="get-transaction-confirmationz-400"></span> 400
Status: Bad Request

###### <span id="get-transaction-confirmationz-400-schema"></span> Schema

### <span id="get-txn-fees"></span> get txn fees (*GetTxnFees*)

```
GET /v1/estimate_txn_fee
```

Estimate transaction fees
Returns an on-chain calculation of the fee based on the provided txn data (in SAS which is the indivisible unit of ZCN coin, 1 ZCN = 10^10 SAS). Txn data is provided in the body of the request.

#### Consumes
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-txn-fees-200) | OK | TxnFeeResponse |  | [schema](#get-txn-fees-200-schema) |

#### Responses


##### <span id="get-txn-fees-200"></span> 200 - TxnFeeResponse
Status: OK

###### <span id="get-txn-fees-200-schema"></span> Schema
   
  

[TxnFeeResponse](#txn-fee-response)

### <span id="get-txn-fees-table"></span> get txn fees table (*GetTxnFeesTable*)

```
GET /v1/fees_table
```

Get transaction fees table
Returns the transaction fees table based on the latest finalized block.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-txn-fees-table-200) | OK | FeesTableResponse |  | [schema](#get-txn-fees-table-200-schema) |

#### Responses


##### <span id="get-txn-fees-table-200"></span> 200 - FeesTableResponse
Status: OK

###### <span id="get-txn-fees-table-200-schema"></span> Schema
   
  

[FeesTableResponse](#fees-table-response)

## Models

### <span id="block"></span> Block


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ChainID | string| `string` |  | |  |  |
| Hash | string| `string` |  | |  |  |
| K | int64 (formatted integer)| `int64` |  | |  |  |
| LatestFinalizedMagicBlockHash | string| `string` |  | |  |  |
| LatestFinalizedMagicBlockRound | int64 (formatted integer)| `int64` |  | |  |  |
| MagicBlockNumber | int64 (formatted integer)| `int64` |  | |  |  |
| MinerID | string| `string` |  | |  |  |
| N | int64 (formatted integer)| `int64` |  | |  |  |
| PrevBlockVerificationTickets | [][VerificationTicket](#verification-ticket)| `[]*VerificationTicket` |  | |  |  |
| PrevHash | string| `string` |  | |  |  |
| PreviousMagicBlockHash | string| `string` |  | |  |  |
| Round | int64 (formatted integer)| `int64` |  | |  |  |
| RoundRandomSeed | int64 (formatted integer)| `int64` |  | |  |  |
| RoundTimeoutCount | int64 (formatted integer)| `int64` |  | |  |  |
| RunningTxnCount | int64 (formatted integer)| `int64` |  | |  |  |
| Signature | string| `string` |  | |  |  |
| StartingRound | int64 (formatted integer)| `int64` |  | |  |  |
| StateChangesCount | int64 (formatted integer)| `int64` |  | | StateChangesCount represents the state changes number in client state of current block.</br>this will be used to verify the state changes acquire from remote |  |
| T | int64 (formatted integer)| `int64` |  | |  |  |
| Txns | [][Transaction](#transaction)| `[]*Transaction` |  | | The entire transaction payload to represent full block |  |
| VerificationTickets | [][VerificationTicket](#verification-ticket)| `[]*VerificationTicket` |  | |  |  |
| Version | string| `string` |  | | Version of the entity |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |
| miners | [Pool](#pool)| `Pool` |  | |  |  |
| mpks | [Mpks](#mpks)| `Mpks` |  | |  |  |
| sharders | [Pool](#pool)| `Pool` |  | |  |  |
| share_or_signs | [GroupSharesOrSigns](#group-shares-or-signs)| `GroupSharesOrSigns` |  | |  |  |
| state_hash | [Key](#key)| `Key` |  | |  |  |



### <span id="block-fee-stats-response"></span> BlockFeeStatsResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| MaxFee | int64 (formatted integer)| `int64` |  | |  |  |
| MeanFee | int64 (formatted integer)| `int64` |  | |  |  |
| MinFee | int64 (formatted integer)| `int64` |  | |  |  |



### <span id="block-response"></span> BlockResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| MerkleTree | []string| `[]string` |  | | Will be returned if only merkle tree is requested. |  |
| block | [Block](#block)| `Block` |  | |  |  |
| header | [BlockSummary](#block-summary)| `BlockSummary` |  | |  |  |



### <span id="block-summary"></span> BlockSummary


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hash | string| `string` |  | |  |  |
| K | int64 (formatted integer)| `int64` |  | |  |  |
| MagicBlockNumber | int64 (formatted integer)| `int64` |  | |  |  |
| MerkleTreeRoot | string| `string` |  | |  |  |
| MinerID | string| `string` |  | |  |  |
| N | int64 (formatted integer)| `int64` |  | |  |  |
| NumTxns | int64 (formatted integer)| `int64` |  | |  |  |
| PreviousMagicBlockHash | string| `string` |  | |  |  |
| ReceiptMerkleTreeRoot | string| `string` |  | |  |  |
| Round | int64 (formatted integer)| `int64` |  | |  |  |
| RoundRandomSeed | int64 (formatted integer)| `int64` |  | |  |  |
| StartingRound | int64 (formatted integer)| `int64` |  | |  |  |
| StateChangesCount | int64 (formatted integer)| `int64` |  | |  |  |
| T | int64 (formatted integer)| `int64` |  | |  |  |
| Version | string| `string` |  | | Version of the entity |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |
| miners | [Pool](#pool)| `Pool` |  | |  |  |
| mpks | [Mpks](#mpks)| `Mpks` |  | |  |  |
| sharders | [Pool](#pool)| `Pool` |  | |  |  |
| share_or_signs | [GroupSharesOrSigns](#group-shares-or-signs)| `GroupSharesOrSigns` |  | |  |  |
| state_hash | [Key](#key)| `Key` |  | |  |  |



### <span id="chain-info"></span> ChainInfo


> swgger:model
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| latest_finalized_block | [BlockSummary](#block-summary)| `BlockSummary` |  | |  |  |



### <span id="chain-stats"></span> ChainStats


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Count | int64 (formatted integer)| `int64` |  | | Number of finalized blocks generated in the block chain since genesis. |  |
| CurrentRound | int64 (formatted integer)| `int64` |  | | The number that represents the current round of the blockchain. |  |
| LastFinalizedRound | int64 (formatted integer)| `int64` |  | | The number that represents the round that generated the latest finalized block. |  |
| Max | double (formatted number)| `float64` |  | | Maximum finalization time of a block, in milliseconds. |  |
| Mean | double (formatted number)| `float64` |  | | Mean (Average) finalization time of a block, in milliseconds. |  |
| Min | double (formatted number)| `float64` |  | | Minimum finalization time of a block, in milliseconds. |  |
| Percentile50 | double (formatted number)| `float64` |  | | The block finalization time value, in milliseconds, which the specified percentage of block finalization events lie below. |  |
| Percentile90 | double (formatted number)| `float64` |  | | The block finalization time value, in milliseconds, which the specified percentage of block finalization events lie below. |  |
| Percentile95 | double (formatted number)| `float64` |  | | The block finalization time value, in milliseconds, which the specified percentage of block finalization events lie below. |  |
| Percentile99 | double (formatted number)| `float64` |  | | The block finalization time value, in milliseconds, which the specified percentage of block finalization events lie below. |  |
| Rate1 | double (formatted number)| `float64` |  | | The moving average rate of occurrence of block finalization events per second during the specified time window. |  |
| Rate15 | double (formatted number)| `float64` |  | | The moving average rate of occurrence of block finalization events per second during the specified time window. |  |
| Rate5 | double (formatted number)| `float64` |  | | The moving average rate of occurrence of block finalization events per second during the specified time window. |  |
| RateMean | double (formatted number)| `float64` |  | | The overall mean rate of occurrence of block finalization events per second. |  |
| RunningTxnCount | int64 (formatted integer)| `int64` |  | | The total count of all transactions included in all the blocks generated by the blockchain. |  |
| StdDev | double (formatted number)| `float64` |  | | Standard deviation of the finalization time of a block from the mean number, in milliseconds. |  |
| delta | [Duration](#duration)| `Duration` |  | |  |  |



### <span id="client"></span> Client


> Client - data structure that holds the client data
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` |  | |  |  |
| PublicKey | string| `string` |  | | The public key of the client |  |
| Version | string| `string` |  | | Version of the entity |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |



### <span id="coin"></span> Coin


> go:generate msgp -io=false -tests=false -v
  



| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| Coin | uint64 (formatted integer)| uint64 | | go:generate msgp -io=false -tests=false -v |  |



### <span id="confirmation"></span> Confirmation


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| BlockHash | string| `string` |  | |  |  |
| Hash | string| `string` |  | |  |  |
| MerkleTreeRoot | string| `string` |  | |  |  |
| MinerID | string| `string` |  | |  |  |
| PreviousBlockHash | string| `string` |  | |  |  |
| ReceiptMerkleTreeRoot | string| `string` |  | |  |  |
| Round | int64 (formatted integer)| `int64` |  | |  |  |
| RoundRandomSeed | int64 (formatted integer)| `int64` |  | |  |  |
| StateChangesCount | int64 (formatted integer)| `int64` |  | |  |  |
| Status | int64 (formatted integer)| `int64` |  | |  |  |
| Version | string| `string` |  | |  |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |
| merkle_tree_path | [MTPath](#m-t-path)| `MTPath` |  | |  |  |
| receipt_merkle_tree_path | [MTPath](#m-t-path)| `MTPath` |  | |  |  |
| txn | [Transaction](#transaction)| `Transaction` |  | |  |  |



### <span id="confirmation-response"></span> ConfirmationResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Error | string| `string` |  | | Error message if any. |  |
| confirmation | [Confirmation](#confirmation)| `Confirmation` |  | |  |  |
| latest_finalized_block | [BlockSummary](#block-summary)| `BlockSummary` |  | |  |  |



### <span id="creation-date-field"></span> CreationDateField


> go:generate msgp -io=false -tests=false -v
CreationDateField - Can be used to add a creation date functionality to an entity */
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |



### <span id="d-k-g-key-share"></span> DKGKeyShare


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` |  | |  |  |
| Message | string| `string` |  | |  |  |
| Share | string| `string` |  | |  |  |
| Sign | string| `string` |  | |  |  |



### <span id="duration"></span> Duration


> A Duration represents the elapsed time between two instants
as an int64 nanosecond count. The representation limits the
largest representable duration to approximately 290 years.
  



| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| Duration | int64 (formatted integer)| int64 | | A Duration represents the elapsed time between two instants</br>as an int64 nanosecond count. The representation limits the</br>largest representable duration to approximately 290 years. |  |



### <span id="explorer-stats"></span> ExplorerStats


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AverageBlockSize | int64 (formatted integer)| `int64` |  | | Average size of the block in terms of number of transaction. |  |
| LastFinalizedRound | int64 (formatted integer)| `int64` |  | | Latest finalized round as known to the sharder. |  |
| MeanScanBlockStatsTime | double (formatted number)| `float64` |  | | Mean sharder health check time. |  |
| PrevInvocationCount | uint64 (formatted integer)| `uint64` |  | | How many times health check was invoked for the sharder. |  |
| PrevInvocationScanTime | string| `string` |  | | How long did it take the previous health check invocation to run, in seconds. |  |
| StateHealth | int64 (formatted integer)| `int64` |  | | Number of missing nodes as seen by the sharder. |  |



### <span id="fees-table-response"></span> FeesTableResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ScFeesTableMap | map of [map[string]int64](#map-string-int64)| `map[string]map[string]int64` |  | |  |  |



### <span id="group-shares-or-signs"></span> GroupSharesOrSigns


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Shares | map of [ShareOrSigns](#share-or-signs)| `map[string]ShareOrSigns` |  | |  |  |



### <span id="hash-id-field"></span> HashIDField


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hash | string| `string` |  | |  |  |



### <span id="health-check-response"></span> HealthCheckResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| BuildTag | string| `string` |  | | Build tag of the image used to deploy the sharder. |  |
| NodeType | string| `string` |  | | Should always be "sharder" |  |
| Uptime | string| `string` |  | | Uptime of the sharder in Nanoseconds. |  |
| chain | [ChainInfo](#chain-info)| `ChainInfo` |  | |  |  |



### <span id="id-field"></span> IDField


> go:generate msgp -io=false -tests=false -v
IDField - Useful to embed this into all the entities and get consistent behavior */
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` |  | |  |  |



### <span id="info"></span> Info


> Info - (informal) info of a node that can be shared with other nodes
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AvgBlockTxns | int64 (formatted integer)| `int64` |  | |  |  |
| BuildTag | string| `string` |  | |  |  |
| StateMissingNodes | int64 (formatted integer)| `int64` |  | |  |  |
| miners_median_network_time | [Duration](#duration)| `Duration` |  | |  |  |



### <span id="info-response"></span> InfoResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ChainInfo | [][ChainInfo](#chain-info)| `[]*ChainInfo` |  | |  |  |
| RoundInfo | [][RoundInfo](#round-info)| `[]*RoundInfo` |  | |  |  |



### <span id="key"></span> Key


> Key - a type for the merkle patricia trie node key
  



[]uint8 (formatted integer)

### <span id="m-p-k"></span> MPK


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` |  | |  |  |
| Mpk | []string| `[]string` |  | |  |  |



### <span id="m-t-path"></span> MTPath


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| LeafIndex | int64 (formatted integer)| `int64` |  | |  |  |
| Nodes | []string| `[]string` |  | |  |  |



### <span id="magic-block"></span> MagicBlock


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hash | string| `string` |  | |  |  |
| K | int64 (formatted integer)| `int64` |  | |  |  |
| MagicBlockNumber | int64 (formatted integer)| `int64` |  | |  |  |
| N | int64 (formatted integer)| `int64` |  | |  |  |
| PreviousMagicBlockHash | string| `string` |  | |  |  |
| StartingRound | int64 (formatted integer)| `int64` |  | |  |  |
| T | int64 (formatted integer)| `int64` |  | |  |  |
| miners | [Pool](#pool)| `Pool` |  | |  |  |
| mpks | [Mpks](#mpks)| `Mpks` |  | |  |  |
| sharders | [Pool](#pool)| `Pool` |  | |  |  |
| share_or_signs | [GroupSharesOrSigns](#group-shares-or-signs)| `GroupSharesOrSigns` |  | |  |  |



### <span id="mpks"></span> Mpks


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Mpks | map of [MPK](#m-p-k)| `map[string]MPK` |  | |  |  |



### <span id="n-o-id-field"></span> NOIDField


> NOIDFied - used when we just want to create a datastore entity that doesn't
have it's own id (like 1-to-many) that is only required to send it around with the parent key */
  



[interface{}](#interface)

### <span id="node"></span> Node


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Description | string| `string` |  | |  |  |
| Host | string| `string` |  | |  |  |
| ID | string| `string` |  | |  |  |
| InPrevMB | boolean| `bool` |  | |  |  |
| N2NHost | string| `string` |  | |  |  |
| Path | string| `string` |  | |  |  |
| Port | int64 (formatted integer)| `int64` |  | |  |  |
| PublicKey | string| `string` |  | | The public key of the client |  |
| SetIndex | int64 (formatted integer)| `int64` |  | |  |  |
| Status | int64 (formatted integer)| `int64` |  | |  |  |
| Version | string| `string` |  | | Version of the entity |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |
| info | [Info](#info)| `Info` |  | |  |  |
| type | [NodeType](#node-type)| `NodeType` |  | |  |  |



### <span id="node-type"></span> NodeType


  

| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| NodeType | int8 (formatted integer)| int8 | |  |  |



### <span id="pool"></span> Pool


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| NodesMap | map of [Node](#node)| `map[string]Node` |  | |  |  |
| type | [NodeType](#node-type)| `NodeType` |  | |  |  |



### <span id="pool-member-info"></span> PoolMemberInfo


> PoolMemberInfo of a pool member
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| N2NHost | string| `string` |  | |  |  |
| Port | string| `string` |  | |  |  |
| PublicKey | string| `string` |  | |  |  |
| Type | string| `string` |  | |  |  |



### <span id="pool-members-info"></span> PoolMembersInfo


> PoolMembersInfo array of pool memebers
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| MembersInfo | [][PoolMemberInfo](#pool-member-info)| `[]*PoolMemberInfo` |  | |  |  |



### <span id="round-info"></span> RoundInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| MultiNotarizedBlocksCount | int8 (formatted integer)| `int8` |  | | count of rounds with multiple notarized blocks. |  |
| NotarizedBlocksCount | int8 (formatted integer)| `int8` |  | |  |  |
| Round | int64 (formatted integer)| `int64` |  | |  |  |
| TimeStamp | string| `string` |  | |  |  |
| ZeroNotarizedBlocksCount | int8 (formatted integer)| `int8` |  | | count of rounds with no notarization for any blocks |  |



### <span id="share-or-signs"></span> ShareOrSigns


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` |  | |  |  |
| ShareOrSigns | map of [DKGKeyShare](#d-k-g-key-share)| `map[string]DKGKeyShare` |  | |  |  |



### <span id="state"></span> State


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Nonce | int64 (formatted integer)| `int64` |  | | Latest nonce used by the client wallet. |  |
| Round | int64 (formatted integer)| `int64` |  | | Latest round when the latest txn happened. |  |
| TxnHash | string| `string` |  | | Latest transaction run by the client wallet. |  |
| balance | [Coin](#coin)| `Coin` |  | |  |  |



### <span id="string-map"></span> StringMap


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Fields | map of string| `map[string]string` |  | |  |  |



### <span id="timestamp"></span> Timestamp


> go:generate msgp -io=false -tests=false -v
Timestamp - just a wrapper to control the json encoding */
  



| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| Timestamp | int64 (formatted integer)| int64 | | go:generate msgp -io=false -tests=false -v</br>Timestamp - just a wrapper to control the json encoding */ |  |



### <span id="transaction"></span> Transaction


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ChainID | string| `string` | ✓ | | ChainID - the chain id of the transaction |  |
| ClientID | string| `string` | ✓ | | ClientID of the client issuing the transaction |  |
| Hash | string| `string` |  | |  |  |
| Nonce | int64 (formatted integer)| `int64` | ✓ | | Nonce - the nonce associated with the transaction |  |
| OutputHash | string| `string` | ✓ | | OutputHash - the hash of the transaction output |  |
| PublicKey | string| `string` | ✓ | | Public key of the client issuing the transaction |  |
| Signature | string| `string` | ✓ | | Signature - Issuer signature of the transaction |  |
| Status | int64 (formatted integer)| `int64` | ✓ | | Status - the status of the transaction |  |
| ToClientID | string| `string` | ✓ | | ToClientID - the client id of the recipient, the other party in the transaction. It can be a client id or the address of a smart contract |  |
| TransactionData | string| `string` | ✓ | | TransactionData - the data associated with the transaction |  |
| TransactionOutput | string| `string` | ✓ | | TransactionOutput - the output of the transaction |  |
| TransactionType | int64 (formatted integer)| `int64` | ✓ | | TransactionType - the type of the transaction. </br>Possible values are:</br>0: TxnTypeSend - A transaction to send tokens to another account, state is maintained by account.</br>10: TxnTypeData - A transaction to just store a piece of data on the block chain.</br>1000: TxnTypeSmartContract - A smart contract transaction type. |  |
| Version | string| `string` |  | | Version of the entity |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` | ✓ | |  |  |
| transaction_fee | [Coin](#coin)| `Coin` | ✓ | |  |  |
| transaction_value | [Coin](#coin)| `Coin` | ✓ | |  |  |



### <span id="txn-fee-response"></span> TxnFeeResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Fee | string| `string` |  | |  |  |



### <span id="unverified-block-body"></span> UnverifiedBlockBody


> UnverifiedBlockBody - used to compute the signature
This is what is used to verify the correctness of the block & the associated signature
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| LatestFinalizedMagicBlockHash | string| `string` |  | |  |  |
| LatestFinalizedMagicBlockRound | int64 (formatted integer)| `int64` |  | |  |  |
| MinerID | string| `string` |  | |  |  |
| PrevBlockVerificationTickets | [][VerificationTicket](#verification-ticket)| `[]*VerificationTicket` |  | |  |  |
| PrevHash | string| `string` |  | |  |  |
| Round | int64 (formatted integer)| `int64` |  | |  |  |
| RoundRandomSeed | int64 (formatted integer)| `int64` |  | |  |  |
| RoundTimeoutCount | int64 (formatted integer)| `int64` |  | |  |  |
| Txns | [][Transaction](#transaction)| `[]*Transaction` |  | | The entire transaction payload to represent full block |  |
| Version | string| `string` |  | | Version of the entity |  |
| creation_date | [Timestamp](#timestamp)| `Timestamp` |  | |  |  |
| state_hash | [Key](#key)| `Key` |  | |  |  |



### <span id="verification-ticket"></span> VerificationTicket


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Signature | string| `string` |  | |  |  |
| VerifierID | string| `string` |  | |  |  |



### <span id="version-field"></span> VersionField


> go:generate msgp -io=false -tests=false -v
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Version | string| `string` |  | | Version of the entity |  |

