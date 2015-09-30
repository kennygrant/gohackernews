/* JS for stories */

DOM.Ready(function() {
  // Watch story form to fetch title of story page
  SetSubmitStoryName();
});


// SetSubmitStoryName sets the story name from a URL
//  attempt to extract a last param from URL, and set name to a munged version of that
function SetSubmitStoryName() {

  DOM.On('.active_url_field', 'change', function(e) {
    var field = DOM.First('.active_name_field')
    if (field.value == "") {
      field.value = urlToSentenceCase(this.value);
    }
  });

}

// Change a URL to a sentence for SetSubmitStoryName
function urlToSentenceCase(url) {
  var parts, name
  url = url.replace(/\/$/, ""); // remove trailing /
  parts = url.split("/"); // now split on /
  name = parts[parts.length - 1]; // last part of string after last /
  name = name.replace(/^\d*-/, ""); // remove prefix numerals with dash (common on id based keys)
  name = name.replace(/\..*$/, ""); // remove .html etc extensions
  name = name.replace(/[_\-+]/g, " "); // remove all - or + or _ in string, replacing with space
  name = name.trim(); // remove whitespace trailing or leading
  name = name.toLowerCase(); // all lower
  name = name[0].toUpperCase() + name.substring(1); // Sentence case
  //  name = name.split(" ").map(function(i){return i[0].toUpperCase() + i.substring(1)}).join(" "); // titlecase
  return name
}