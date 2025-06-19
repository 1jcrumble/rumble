# Rumble Advertising Center Golang Challenge

This challenge is running in a simulated environment and is using a simple one file approach to understand how you think about and solve problems. The boundaries are that the results need to be accessible via curl, the code needs to compile and run within this environment, and the instructions to get the results should be clear. 

## Objectives

This challenge has two components: written and code. For the written, create a well formatted markdown document outlining your responses. Point form and succinct responses are valued.

### Written
- Review the existing code, what is the work that needs to be done if you were to take this code into production in your opinion? How does this mock up / challenge exercise differ from what you would expect to see in the 'real world'? 
- Also use this document to highlight what you have done in the code (and why)

## Code
Please complete the following (based on and using the code provided)
- Right now the project uses mock data, please externalize this so that the program can injest both a JSON file and an external url as a source of mock transactions
- Add a command line option to be able to specify the data source file name, e.g: 
```
# json file
./main --transactions transactions.json

# external url
./main --transactions https://domain.com/transactions.json
```
- Instead of displaying the PAN with GetTransactions it is preferred to only display the last four digits and replace the rest of the PAN with `*`; create a function to achieve this and ensure that all output is handled in this way
- Create an endpoint that returns the transactions ordered by descending posted_timestamp 
- Create a test for GetTransactions and your new functions

## Review of Existing Code

### Gaps Before Production
- **No Configuration/Flags**: Lacks CLI flags or config files to specify input sources or port.
- **Hardcoded Port**: Uses fixed port `:8000` instead of being configurable.
- **PAN Security**: Full PANs are printed in logs/output, which is a major PCI compliance issue.
- **No Error Logging**: HTTP handlers do not log errors.
- **No Tests**: No unit or integration tests exist.
- **No Input Validation**: PAN and timestamp formats are not validated.
- **No Timestamps Parsed**: `PostedTimeStamp` is only a string, not parsed into a time object.

### Differences From Real World
- In production, we’d have:
    - Configurable ports and paths
    - Authentication and access control
    - Observability: logging, metrics, tracing
    - Input/output validation and schema checks
    - Structured project layout with modular code

## What I Implemented

### Features
- Added `--transactions` CLI flag to specify a local file or remote URL
- Implemented external transaction loading from JSON file or HTTP(S)
- Masked PANs to show only last 4 digits
- Sorted transactions by `posted_timestamp` in descending order
- Added test coverage for:
    - PAN masking
    - Sorting
    - `/transactions` endpoint

### Why
- **CLI Flag**: Enables flexibility in data sourcing.
- **PAN Masking**: Protects sensitive data.
- **Sorting**: Common requirement for transaction UIs or exports.
- **Tests**: To ensure correctness of sensitive logic.