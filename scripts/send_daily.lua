--cron: 20 7,19 * * *

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

local messages = {
    daily = "Bom dia @here! Bora compartilhar sua Daily? 🐶🤖",
    alert_marco = string.format("Hello, <@%s>! Me respeita soooooo celoooko??", user_map.marco)
}

if not discord.is_work_day() then
    print("it's weekend! 🏖️🍻")
    discord.send_message(channel_map.geral, "É fim de semana galera, lembrem de descansar e aproveitar! 🏖️🍻")
    return
end

if discord.send_message(channel_map.daily, messages.daily) then
    -- check if it's not a weekend
    print("message sent! 🎉")
end
