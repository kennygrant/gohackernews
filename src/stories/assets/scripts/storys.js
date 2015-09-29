/* JS for stories */

DOM.Ready(function() {
  // Watch story form to fetch title of story page
  SetSubmitStoryName();

});


function SetSubmitStoryName() {

  // From the url, attempt to extract a last param, and set name to a munged version of that
  DOM.On('.active_url_field', 'change', function(e) {
    DOM.First('.active_name_field').value = urlToSentenceCase(this.value);
  });

}


function urlToSentenceCase(url) {
  var parts, name
  url = url.replace(/\/$/, ""); // remove trailing /
  parts = url.split("/"); // now split on /
  name = parts[parts.length - 1]; // last part of string after last /
  name = name.replace(/^\d*-/, ""); // remove prefix numerals with dash (common on id based keys)
  name = name.replace(/\..*$/, ""); // remove .html etc extensions
  name = name.replace(/-/g, " "); // remove all - in string, replacing with space
  name = name.trim(); // remove whitespace trailing or leading
  name = name.toLowerCase(); // all lower
  name = name[0].toUpperCase() + name.substring(1); // Sentence case
  //  name = name.split(" ").map(function(i){return i[0].toUpperCase() + i.substring(1)}).join(" "); // titlecase
  return name
}