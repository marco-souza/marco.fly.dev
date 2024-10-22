--cron: 20 4,16 * * *
if not discord.is_work_day() then
    print("it's weekend! 🏖️🍻")
    -- discord.send_message(channel_map.geral, "É fim de semana galera, lembrem de descansar e aproveitar! 🏖️🍻")
    return
end

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
    dailys = {
        "Bom dia @here! Bora compartilhar sua Daily? 🐶🤖",
        "Alou @here! Quem vai compartilhar sua Daily hoje? 🐶🤖",
        "Galeralera @here! Bora compartilhar sua Daily? 🐶🤖",
        "@here Aobo bão? Quem vai participar da Daily hoje? 🐶🤖",
        "Pulei no muro e cai de costa, @here quem vai postar a daily hoje seus bosta? 🐶🤖",
        string.format(
            "@here quem não postar a daily hoje vai ter %d comentários no seu próximo code review 🐶🤖",
            math.random(13, 42)
        ),
    },
    alert_marco = string.format("Hello, <@%s>! Me respeita soooooo celoooko??", user_map.marco),
}

local msg = messages.dailys[math.random(#messages.dailys)]
if discord.send_message(channel_map.daily, msg) then
    -- check if it's not a weekend
    print("message sent! 🎉")
end
