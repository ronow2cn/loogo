#!/bin/env lua

-- ============================================================================

local DIR_PROTO   = "src/proto"

-- ============================================================================

local T = {
    -- ["c_s"] = {
    --      files   = {...},
    --      structs = {{key = "c", msgid = 1, tp = xxx}, ...},
    -- }
}

-- ============================================================================

local function throw(fmt, ...)
    error(string.format(fmt, ...))
end

local function execute(fmt, ...)
    os.execute(string.format(fmt, ...))
end

local function read_file(fn)
    local f = io.open(fn, "rt")
    if not f then throw("open file failed: %s", fn) end

    local text = f:read("*a")
    f:close()

    return text
end

local function save_file(fn, text)
    local f = io.open(fn, "wt")
    if not f then throw("save file failed: %s", fn) end

    f:write(text)
    f:close()
end

local printer_arr

local function print_start()
    printer_arr = {}
end

local function print_end()
    return table.concat(printer_arr, "\n")
end

local function printf(indent, fmt, ...)
    if indent then
        table.insert(printer_arr, string.format("%s%s", string.rep(" ", indent * 4), string.format(fmt, ...)))
    else
        table.insert(printer_arr, "")
    end
end

local function file_exists(fn)
    local f = io.open(fn, "r")
    if f then
        f:close()
        return true
    else
        return false
    end
end

-- ----------------------------------------------------------------------------

local function print_msg_interface()
    printf(0, "package msg")
    printf()
    printf(0, "import (")
    printf(0, '    "gopkg.in/vmihailenco/msgpack.v2"')
    printf(0, ")")
    printf()
    printf(0, "type Message interface {")
    printf(0, "    MsgId() uint32")
    printf(0, "}")
    printf()
    printf(0, "func Marshal(m Message) ([]byte, error) {")
    printf(0, "    return msgpack.Marshal(m)")
    printf(0, "}")
    printf()
    printf(0, "func Unmarshal(b []byte, obj interface{}) error {")
    printf(0, "    return msgpack.Unmarshal(b, obj)")
    printf(0, "}")
    printf()
end

-- ----------------------------------------------------------------------------

local function print_struct_creator(cnn_arr)
    printf(0, "var MsgCreators = map[uint32]func() Message{")

    for _, cnn in ipairs(cnn_arr) do
        local prefix = cnn[1]

        for _, struct in ipairs(T[prefix].structs) do
            printf(1, "%d: func() Message {", struct.msgid)
            printf(1, "    return &%s{}", struct.tp)
            printf(1, "},")
        end
    end

    printf(0, "}")
    printf()
end

local function print_struct_msgid(cnn_arr)
    for _, cnn in ipairs(cnn_arr) do
        local prefix = cnn[1]

        for _, struct in ipairs(T[prefix].structs) do
            printf(0, "func (self *%s) MsgId() uint32 {", struct.tp)
            printf(0, "    return %d", struct.msgid)
            printf(0, "}")
            printf()
        end
    end
end

-- ----------------------------------------------------------------------------

local function print_handler_dt()
    printf(0, "package msg")
    printf()
    printf(0, "var MsgHandlers = map[uint32]func(message Message, ctx interface{}){}")
    printf()
    printf(0, "func Handler(msgid uint32, h func(message Message, ctx interface{})) {")
    printf(0, "    MsgHandlers[msgid] = h")
    printf(0, "}")
    printf()
end

local function print_handler_map(svr, cnn_arr)
    -- sort imports
    local imports = {}
    for _, cnn in ipairs(cnn_arr) do
        local prefix = cnn[1]
        table.insert(imports, string.format("%s/handler/%s", svr, prefix))
    end
    table.insert(imports, string.format("%s/msg", svr))
    table.sort(imports)

    -- print
    printf(0, "package handler")
    printf()
    printf(0, "import (")

    for _, v in ipairs(imports) do
        printf(0, "    %q", v)
    end

    printf(0, ")")
    printf()
    printf(0, "func Init() {")

    for _, cnn in ipairs(cnn_arr) do
        local prefix = cnn[1]
        local key    = string.lower(cnn[2])

        for _, struct in ipairs(T[prefix].structs) do
            if struct.key == key then
                printf(1, "msg.Handler(%d, %s.%s)", struct.msgid, prefix, struct.tp)
            end
        end
    end
    printf(0, "}")
    printf()
end

-- ----------------------------------------------------------------------------

local function print_handler(svr, prefix, tp)
    printf(0, "package %s", prefix)
    printf()
    printf(0, "import (")
    printf(0, '    "%s/msg"', svr)
    printf(0, ")")
    printf()
    printf(0, "func %s(message msg.Message, ctx interface{}) {", tp)
    printf(0, "    req := message.(*msg.%s)", tp)
    printf(0, "    req = req")
    printf(0, "}")
    printf()
end

-- ============================================================================

local function prepare()
    local f = io.popen(string.format("find %s -maxdepth 1 -type f -name '*.go'|sort", DIR_PROTO))
    if not f then throw("finding proto files failed") end

    for fn in f:lines() do
        local prefix = string.match(fn, "^"..DIR_PROTO.."/(%w+_%w+)%.")
        if not prefix then throw("invalid proto file prefix: %s", fn) end

        local t2 = T[prefix]
        if not t2 then
            t2 = {files = {}, structs = {}}
            T[prefix] = t2
        end

        -- append file
        table.insert(t2.files, fn)

        -- append struct
        local text = read_file(fn)
        for tp, msgid in string.gmatch(text, "type%s+([%w_]+)%s+struct%s*{%s*//%s*msgid:%s*(%d+)") do
            local key = string.match(tp, "^(%w+)_")
            if not key then throw("invalid key for struct: %s", tp) end

            key = string.lower(key)
            table.insert(t2.structs, {key = key, msgid = msgid, tp = tp})
        end
    end
    f:close()
end

local function gen_msg_file(svr, cnn_arr)
    local outdir = string.format("src/%s/msg", svr)

    execute("rm -rf %s",   outdir)
    execute("mkdir -p %s", outdir)

    for _, cnn in ipairs(cnn_arr) do
       local prefix = cnn[1]

       for _, fn in ipairs(T[prefix].files) do
           execute(string.format("cp %s %s", fn, outdir))
       end
   end
end

local function gen_msg_interface(svr, cnn_arr)
    local outfile, text

    -- print
    print_start()
    print_msg_interface()
    text = print_end()

    -- save
    outfile = string.format("src/%s/msg/message.go", svr)
    save_file(outfile, text)
end

local function gen_struct_map(svr, cnn_arr)
    local outfile, text

    -- print
    print_start()
    printf(0, "package msg")
    printf()
    print_struct_creator(cnn_arr)
    print_struct_msgid(cnn_arr)
    text = print_end()

    -- save
    outfile = string.format("src/%s/msg/idstruct.go", svr)
    save_file(outfile, text)
end

local function gen_handler_map(svr, cnn_arr)
    local outfile, text

    -- handler dt
    print_start()
    print_handler_dt()
    text = print_end()

    -- save
    outfile = string.format("src/%s/msg/handler.go", svr)
    save_file(outfile, text)

    -- handler map
    print_start()
    print_handler_map(svr, cnn_arr)
    text = print_end()

    -- save
    execute("mkdir -p src/%s/handler", svr)
    outfile = string.format("src/%s/handler/init.go", svr)
    save_file(outfile, text)
end

local function gen_handlers(svr, cnn_arr)
    local outfile, text

    for _, cnn in ipairs(cnn_arr) do
        local prefix = cnn[1]
        local key    = cnn[2]

        execute("mkdir -p src/%s/handler/%s", svr, prefix)

        for _, struct in ipairs(T[prefix].structs) do
            if struct.key == key then
                outfile = string.format("src/%s/handler/%s/%s.go", svr, prefix, struct.tp)
                if not file_exists(outfile) then
                    -- print
                    print_start()
                    print_handler(svr, prefix, struct.tp)
                    text = print_end()

                    -- save
                    save_file(outfile, text)
                end
            end
        end
    end
end

-- ============================================================================

local function gen_proto(svr, cnn_arr)
    print("generating proto for "..svr)

    gen_msg_file     (svr, cnn_arr)
    gen_msg_interface(svr, cnn_arr)
    gen_struct_map   (svr, cnn_arr)
    gen_handler_map  (svr, cnn_arr)
    gen_handlers     (svr, cnn_arr)
end

-- ============================================================================
--                            Main
-- ============================================================================

prepare()

print()

gen_proto("gate", {
    {"c_gw",  "c" },
    {"gw_gs", "gs"},
})

gen_proto("game", {
    {"gw_gs", "gw"},
    {"c_gs",  "c" },
})

gen_proto("client", {
    {"c_gw",  "gw"},
    {"c_gs",  "gs"},
})

print()
