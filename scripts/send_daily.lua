local channel_map = {
    daily = "742025222126960760",
    testing = "1243281325780107444",
}

if not discord.is_work_day() then
    return
end

message_sent = discord.send_message(
    channel_map.daily, "Hello, World!"
)
if message_sent then
    -- check if it's not a weekend
    print("message sent! 🎉")
end
