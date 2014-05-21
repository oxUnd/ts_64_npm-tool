<!DOCTYPE html>
{{framework "static/lib/mod.js"}}
<html>
<head>
    <title>{{.title}}</title>
    {{require "static/base.css"}}
    {{require "static/lib/jquery.js"}}
</head>
<body>
<div class="layout">
    {{template "widget/header/header.tpl"}}

    <div id="container">
        <form method="post" action="/new.do">
            <legend>提交FIS插件；如 fis@1.7.9 </legend>
            <label><input id="keyword_search" type="text" name="plg" /></label>
            <input type="submit" />
        </form>
        <div id="message">
            <div id="loading" style="display:none">
            <div id="fountainG">
                <div id="fountainG_1" class="fountainG">
                </div>
                <div id="fountainG_2" class="fountainG">
                </div>
                <div id="fountainG_3" class="fountainG">
                </div>
                <div id="fountainG_4" class="fountainG">
                </div>
                <div id="fountainG_5" class="fountainG">
                </div>
                <div id="fountainG_6" class="fountainG">
                </div>
                <div id="fountainG_7" class="fountainG">
                </div>
                <div id="fountainG_8" class="fountainG">
                </div>
            </div>
            </div>
        </div>
        <table id="list-table">
            <thead>
                <tr>
                    <th>插件({{.components|len}})</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tr>
                <td></td>
                <td></td>
                <td>
                    <button typ="refresh">refresh</button>
                </td>
            </tr>
            <tbody>
                {{range .components}}
                    <tr>
                        <td><a href="https://www.npmjs.org/package/{{.name}}" target="_blank">{{.name}}</a>@{{.version}}</td>
                        <td>
                        {{if eq .status 0}} 已安装
                        {{else if eq .status 1}} 未安装
                        {{else if eq .status 2}} 安装中
                        {{end}}
                        </td>
                        <td>
                            <button class="install_btn" comp="{{.name}}@{{.version}}" typ="install">install</button>
                        </td>
                    </tr>
                {{end}}
            </tbody>
        </table>    
    </div>

    <script>
        $(document.getElementById('list-table')).click(function (e) {
            if (e.target && e.target.nodeName == 'BUTTON') {
                var el = $(e.target);
                el.attr('disabled', 'disabled');
                var type = el.attr('typ');
                var comp = el.attr('comp');
                $('#loading').show();
                switch (type) {
                    case 'install':
                        $.post('/action.do', {
                            type: type,
                            comp: comp
                        }, function (res) {
                            var s = '';
                            if (res.code == '0') {
                                s = "<h3>安装成功</h3>";
                            } else {
                                s = "<h3>安装失败</h3>";
                            }
                            $('#message').append(s + res.msg.replace(/\n/g, '<br />'));
                            $('#loading').hide()
                        });
                        break;
                    case 'refresh':
                        $.post('/refresh.do', function (result) {
                            $('#loading').hide()
                        });

                        break;
                }
            } 
        });

        $("#keyword_search").keyup(function(){

            var term=$(this).val()
            if( term != "")
            {

                $("#list-table tbody>tr").hide();
                $("#list-table td").filter(function(){
                       return $(this).text().toLowerCase().indexOf(term ) >-1
                }).parent("tr").show();
            }
            else
            {
                $("#list-table tbody>tr").show();
            }
        });
    </script>
    
    {{template "widget/footer/footer.tpl"}}

    {{require "page/index.tpl"}}

</div>

</body>
</html>