# Database connection

## ToDo
- [ ] Pass context for connection
- [ ] sometimes workers open two connections during start:
```bash
1:33AM INF Connect component=db data={"database":"invest_platform","host":"127.0.0.1","pid":858,"port":5432,"time":5564173} severity=INFO
1:33AM INF Connect component=db data={"database":"invest_platform","host":"127.0.0.1","pid":858,"port":5432,"time":5564273} severity=INFO
```
