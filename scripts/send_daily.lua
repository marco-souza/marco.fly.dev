local channel_map = {
    daily = "742025222126960760",
    geral = "694542727358185502",
    testing = "1243281325780107444",
}

local user_map = {
    marco = "488746421944582154",
    nicolas = "694959368827043880",
    vitor = "311246917248221184",
    rafael = "139083790344650753",
}

if not discord.is_work_day() then
    return
end

local message = string.format("Hello, <@%s>! Bora acordar cachorro 🐶", user_map.marco)
if discord.send_message(channel_map.testing, message) then
    -- check if it's not a weekend
    print("message sent! 🎉")
end
