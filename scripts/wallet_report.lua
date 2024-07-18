--skip:cron: */10 * * * *  -> doesnt work for US ip

local report = binance.wallet_report()
if telegram.send_message(report) then
    print("wallet report sent! 🎉 ")
end
