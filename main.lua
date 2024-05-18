--@param count number of times to run the function
function times(count, func)
    for i = 1, count do
        print("running function: " .. i)
        func()
    end
end

times(5, function()
    print "Hello, World!"
end)
