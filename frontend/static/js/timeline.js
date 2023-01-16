// Adds padding to one-digit numbers to make them two
function zeroPadded(val) {
    if (val >= 10)
        return val;
    else
        return '0' + val;
}

function pauseTimeline(icon) {
    withTime = false
    playForward = false
    icon.removeClass('fa-pause')
    icon.addClass('fa-play')
    $('#resume').prop('disabled', false)
    $('#time-controller').prop('disabled', false)
    playDistance = 0
}

function playTimeline(icon) {
    withTime = false
    playForward = true
    icon.removeClass('fa-play')
    icon.addClass('fa-pause')
    $('#resume').prop('disabled', true)
    $('#time-controller').prop('disabled', true)
}

// Conveniece functions for dragging cursor
function timelineIsPaused(){
    return $('#play-pause').children().hasClass('fa-play')
}

function timelineIsPlaying(){
    return $('#play-pause').children().hasClass('fa-pause')
}

// Handles timeline components
$(document).ready(function() {
    // Handle timeline buttons
    $('#play-pause').click(function() {
        icon = $(this).children()
        if (icon.hasClass('fa-pause')) {
            pauseTimeline(icon)
        } else if (icon.hasClass('fa-play')) {
            playTimeline(icon)
        }
    });

    $('#resume').click(function() {
        withTime = true
        icon = $('#play-pause').children()
        icon.removeClass('fa-play')
        icon.addClass('fa-pause')
        $('#resume').prop('disabled', true)
        $('#time-controller').prop('disabled', true)
    });

    $('#loginModal').modal('show')

    $('#time-controller').change(function() {

        newTime = new Date($('#time-controller').val())
        timeline.setCustomTime(newTime, cursorID)
        timeline.moveTo(newTime)
        handleTime()
    });
});