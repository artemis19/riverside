$(document).ready(function() {

    $('#login-form').submit(function(e) {

        e.preventDefault();
        username = $('#login-username').val()
        password = $('#login-password').val()
        server = $('#login-server').val()
        port = $('#login-port').val()


        url = `http://${server}:${port}/login`
        $.ajax(url, {
            type: 'POST',
            data: JSON.stringify({
                "username": username,
                "password": password,
            }),
            contentType: 'application/json',

        }, function() {
            // callback function
        }).done(function(data) {
            $('#login-alert').removeClass("alert-light alert-danger alert-success")
            $('#login-alert').addClass("alert-success")
            $('#login-alert').text("Successfully logged in!")
            $('#login-username').prop("disabled", "true")
            $('#login-password').prop("disabled", "true")
            $('#login-cancel').prop("disabled", "true")
            $('#login-button').prop("disabled", "true")

            setTimeout(function() { $('#login-modal').modal('hide'); }, 1000);
        }).fail(function(data) {
            $('#login-alert').removeClass("alert-light alert-danger alert-success")
            $('#login-alert').addClass("alert-danger")
            if (data.responseJSON == undefined) {
                $('#login-alert').text("Failed to contact the server specified. Check your IP and port!")
            } else {
                $('login-alert').text(data.responseJSON.error)
            }
        })
    })


    $('#register-form').submit(function(e) {

        e.preventDefault();
        username = $('#register-username').val()
        password = $('#register-password').val()
        server = $('#register-server').val()
        port = $('#register-port').val()


        url = `http://${server}:${port}/register`
        $.ajax(url, {
            type: 'POST',
            data: JSON.stringify({
                "username": username,
                "password": password,
            }),
            contentType: 'application/json',

        }, function() {
            // callback function
        }).done(function(data) {
            $('#register-alert').removeClass("alert-light alert-danger alert-success")
            $('#register-alert').addClass("alert-success")
            $('#register-alert').text("Successfully registered!")
            $('#register-username').prop("disabled", "true")
            $('#register-password').prop("disabled", "true")
            $('#register-cancel').prop("disabled", "true")
            $('#register-button').prop("disabled", "true")

            setTimeout(function() { $('#register-modal').modal('hide'); }, 1000);
        }).fail(function(data) {
            $('#register-alert').removeClass("alert-light alert-danger alert-success")
            $('#register-alert').addClass("alert-danger")
            if (data.responseJSON == undefined) {
                $('#register-alert').text("Failed to contact the server specified. Check your IP and port!")
            } else {
                $('register-alert').text(data.responseJSON.error)
            }
        })
    })
})