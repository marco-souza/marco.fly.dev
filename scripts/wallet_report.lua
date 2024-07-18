--cron: */10 * * * *

local report = binance.wallet_report()
if telegram.send_message(report) then
    print("wallet report sent! 🎉 ")
end
