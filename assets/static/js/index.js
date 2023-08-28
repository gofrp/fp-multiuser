var $ = layui.$;
$(function () {
    var apiType = {
        Remove: 1,
        Enable: 2,
        Disable: 3
    }

    /**
     * load i18n language
     * @param lang {{}}
     */
    function langLoaded(lang) {
        layui.table.render({
            elem: '#tokenTable',
            url: '/tokens',
            method: 'get',
            where: {},
            dataType: 'json',
            editTrigger: 'dblclick',
            page: navigator.language.indexOf("zh") === 0,
            toolbar: '#toolbarTemplate',
            defaultToolbar: false,
            text: {none: lang['EmptyData']},
            cols: [[
                {type: 'checkbox'},
                {field: 'user', title: lang['User'], width: 150, sort: true},
                {field: 'token', title: lang['Token'], width: 150, sort: true, edit: true},
                {field: 'comment', title: lang['Notes'], sort: true, edit: 'textarea'},
                {field: 'ports', title: lang['AllowedPorts'], sort: true, edit: 'textarea'},
                {field: 'domains', title: lang['AllowedDomains'], sort: true, edit: 'textarea'},
                {field: 'subdomains', title: lang['AllowedSubdomains'], sort: true, edit: 'textarea'},
                {
                    field: 'status',
                    title: lang['Status'],
                    width: 100,
                    templet: '<span>{{d.status? "' + lang['Enable'] + '":"' + lang['Disable'] + '"}}</span>',
                    sort: true
                },
                {title: lang['Operation'], width: 150, toolbar: '#operationTemplate'}
            ]]
        });

        layui.table.on('edit(tokenTable)', function (obj) {
            var field = obj.field;
            var value = obj.value;
            var oldValue = obj.oldValue;
            var before = $.extend(true, {}, obj.data);
            var after = $.extend(true, {}, obj.data);
            if (field === 'token') {
                if (value.trim() === '') {
                    layui.layer.msg(lang['TokenEmpty'])
                    obj.reedit();
                    return;
                }

                before.token = oldValue;
                after.token = value;
            } else if (field === 'comment') {
                before.comment = oldValue;
                after.comment = value;
            } else if (field === 'ports') {
                before.ports = oldValue;
                after.ports = value;
            } else if (field === 'domains') {
                before.domains = oldValue;
                after.domains = value;
            } else if (field === 'subdomains') {
                before.subdomains = oldValue;
                after.subdomains = value;
            }

            update(before, after);
        });

        layui.table.on('toolbar(tokenTable)', function (obj) {
            var id = obj.config.id;
            var checkStatus = layui.table.checkStatus(id);
            switch (obj.event) {
                case 'add':
                    addPopup();
                    break
                case 'remove':
                    batchRemovePopup(checkStatus.data);
                    break
                case 'disable':
                    batchDisablePopup(checkStatus.data);
                    break
                case 'enable':
                    batchEnablePopup(checkStatus.data);
                    break
            }
        });
        layui.table.on('tool(tokenTable)', function (obj) {
            var data = obj.data;
            switch (obj.event) {
                case 'remove':
                    removePopup(data);
                    break;
                case 'disable':
                    disablePopup(data);
                    break;
                case 'enable':
                    enablePopup(data);
                    break
            }
        });

        /**
         * add user popup
         */
        function addPopup() {
            layui.layer.open({
                type: 1,
                title: lang['NewUser'],
                area: ['500px'],
                content: layui.laytpl(document.getElementById('addTemplate').innerHTML).render(),
                btn: [lang['Confirm'], lang['Cancel']],
                btn1: function (index) {
                    if (layui.form.validate('#addUserForm')) {
                        add(layui.form.val('addUserForm'), index)
                    }
                },
                btn2: function (index) {
                    layui.layer.close(index);
                }
            });
        }

        /**
         * add user action
         * @param data user data
         * @param index popup index
         */
        function add(data, index) {
            var loading = layui.layer.load();
            $.ajax({
                url: '/add',
                type: 'post',
                contentType: 'application/json',
                data: JSON.stringify(data),
                success: function (result) {
                    if (result.success) {
                        reloadTable();
                        layui.layer.close(index);
                        layui.layer.msg(lang['OperateSuccess'], function (index) {
                            layui.layer.close(index);
                        });
                    } else {
                        errorMsg(result);
                    }
                },
                complete: function () {
                    layui.layer.close(loading);
                }
            });
        }

        /**
         * update user action
         * @param before data before update
         * @param after data after update
         */
        function update(before, after) {
            var loading = layui.layer.load();
            $.ajax({
                url: '/update',
                type: 'post',
                contentType: 'application/json',
                data: JSON.stringify({
                    before: before,
                    after: after,
                }),
                success: function (result) {
                    if (result.success) {
                        layui.layer.msg(lang['OperateSuccess']);
                    } else {
                        errorMsg(result);
                    }
                },
                complete: function () {
                    layui.layer.close(loading);
                }
            });
        }

        /**
         * batch remove user popup
         * @param data user data list
         */
        function batchRemovePopup(data) {
            if (data.length === 0) {
                layui.layer.msg(lang['ShouldCheckUser']);
                return;
            }
            layui.layer.confirm(lang['ConfirmRemoveUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Remove, data, index);
            });
        }

        /**
         * batch disable user popup
         * @param data user data list
         */
        function batchDisablePopup(data) {
            if (data.length === 0) {
                layui.layer.msg(lang['ShouldCheckUser']);
                return;
            }
            layui.layer.confirm(lang['ConfirmDisableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Disable, data, index);
            });
        }

        /**
         * batch enable user popup
         * @param data user data list
         */
        function batchEnablePopup(data) {
            if (data.length === 0) {
                layui.layer.msg(lang['ShouldCheckUser']);
                return;
            }
            layui.layer.confirm(lang['ConfirmEnableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Enable, data, index);
            });
        }

        /**
         * remove one user popup
         * @param data user data
         */
        function removePopup(data) {
            layui.layer.confirm(lang['ConfirmRemoveUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Remove, [data], index);
            });
        }

        /**
         * disable one user popup
         * @param data user data list
         */
        function disablePopup(data) {
            layui.layer.confirm(lang['ConfirmDisableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Disable, [data], index);
            });
        }

        /**
         * enable one user popup
         * @param data user data list
         */
        function enablePopup(data) {
            layui.layer.confirm(lang['ConfirmEnableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Enable, [data], index);
            });
        }

        /**
         * operate actions
         * @param type {apiType} action type
         * @param data user data list
         * @param index popup index
         */
        function operate(type, data, index) {
            var url;
            if (type === apiType.Remove) {
                url = "/remove";
            } else if (type === apiType.Disable) {
                url = "/disable";
            } else if (type === apiType.Enable) {
                url = "/enable";
            } else {
                layer.layer.msg(lang['OperateError']);
                return;
            }
            var loading = layui.layer.load();
            $.post({
                url: url,
                type: 'post',
                contentType: 'application/json',
                data: JSON.stringify({
                    users: data
                }),
                success: function (result) {
                    if (result.success) {
                        reloadTable();
                        layui.layer.close(index);
                        layui.layer.msg(lang['OperateSuccess'], function (index) {
                            layui.layer.close(index);
                        });
                    } else {
                        errorMsg(result);
                    }
                },
                complete: function () {
                    layui.layer.close(loading);
                }
            });
        }

        /**
         * reload user table
         */
        function reloadTable() {
            var searchData = layui.form.val('searchForm');
            layui.table.reloadData('tokenTable', {
                where: searchData
            }, true)
        }

        /**
         * show error message popup
         * @param result
         */
        function errorMsg(result) {
            var reason = lang['Other Error'];
            if (result.code === 1)
                reason = lang['ParamError'];
            else if (result.code === 2)
                reason = lang['UserExist'];
            layui.layer.msg(lang['OperateFailed'] + ',' + reason)
        }

        /**
         * click event
         */
        $(document).on('click.search', '#searchBtn', function () {
            reloadTable();
            return false;
        }).on('click.reset', '#resetBtn', function () {
            $('#searchForm')[0].reset();
            reloadTable();
            return false;
        });
    }

    var langLoading = layui.layer.load()
    $.getJSON('/lang').done(langLoaded).always(function () {
        layui.layer.close(langLoading);
    });
});