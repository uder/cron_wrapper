# Report requirements

## Begining report
1. Always write to log in full form
2. Add to chat message only if EnableBegin is set
3. Contains basic information about the command
    a. hostname
    b. begin timestamp. ISO8601
    c. pid
    d. command line
    e. unique run id
    f. example
        ```
        2026-04-12T00:00:02 BEGIN "/www/sp/command cache BuildBrandsFromPurchase" 5ca816ae
        host.example.com
        ---
        ```

## Ending report
1. Always placed in log in full form
2. Send to chat
    a. if command exit code is non-zero then Send full form
    b. if EnableStdoutOnSuccess then send full form to chat on success
    c. If DisableChat then don't send to chat
    d. Else send short form to chat
3. Contains information about command execution results. Some of them are taken from beginning report
    a. Start timestamp. ISO8601
    b. exit code. Command exit code
    c. duration
    d. command line
    e. unique run id
    f. hostname
    g. pid
    h. stdout. Command stdout
    i. stderr. Command stderr
4. Output formatting
    a. Result report contains three blocks
        - meta. Contains data about command, process etc
            - Ending ts. ISO8601
            - Process state. INFO if exit code is zero. ERROR if exit code is not zero
            - Exit code
            - duration
            - command line
            - unique run id
            - hostname
            - begining ts. ISO8601
            - process name (?)
            - process pid
            - STDOUT. Command stdout. Omit if empty
            - STDERR. Command stderr. Omit if empty
            - example:
            ```
            2026-04-12T00:00:17 INFO 0 14.27 "/www/sp/command cache SaleAutoWeekHits --days=5" e71e1036
            host.example.com
            begin_at=2026-04-12T00:00:02 process=sh pid=3895691
            ```
        - STDOUT. Contains stdout. Omit if empty
        ```
        STDOUT:
        Found 1792 goods
        ```
        - STDERR. Contains stderr. text. Omit if empty
        ```
        STDERR:
        some stderr
        ```
    b. Full form. Contains all blocks
        - example
        ```
        2026-04-12T00:00:17 INFO 0 14.27 "/www/sp/command cache SaleAutoWeekHits --days=5" e71e1036
        host.example.com
        begin_at=2026-04-12T00:00:02 process=sh pid=3895691
        STDOUT:
        Found 1792 goods
        STDERR:
        some stderr
        ```
    c. Short form. Contains meta only. It is useful for commands that usually don't fail and you still need notification about them
