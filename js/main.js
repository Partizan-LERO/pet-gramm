$(function () {
    $('#word').keyup(function () {
        if ($(this).val().trim().length === 0) {
            $('.search').addClass('disabled');
        } else {
            $('.search').removeClass('disabled');
        }
    });

    $('.search').click(function () {
        if ($(this).hasClass('disabled')) {
            return false;
        }
        let word = $('#word').val();

        axios.get('/word/' + word)
            .then(function (response) {
                if (response.data === null) {
                    $('#anagrams-list').html("I can't find anagrams for this word");
                    return false;
                }

                let ul = "<ul>";

                response.data.filter((item) => {
                    ul += "<li>" + item + "</li>"
                });

                ul += "</ul>";

                $('#anagrams-list').html(ul);
            })
            .catch(function () {
                $('#anagrams-list').html("Hmmm, I can't see the word here!");
            });

        return false;
    });
});