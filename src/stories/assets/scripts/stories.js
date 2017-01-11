/* JS for stories */

DOM.Ready(function() {
  // Watch story form to fetch title of story page
  SetSubmitStoryName();
});


// SetSubmitStoryName sets the story name from a URL
//  attempt to extract a last param from URL, and set name to a munged version of that
function SetSubmitStoryName() {

  DOM.On('.active_url_field', 'change', function(e) {
    var name = DOM.First('.name_field')
    var summary = DOM.First('.summary_field')
    var original_name = name.value
    var original_url = this.value
    
    // First locally fill in the name field 
    if (original_name== "") {
      
      // For github urls, try to fetch some more info 
      if (original_url.startsWith('https://github.com/')) {
      
        var url = original_url.replace('https://github.com/','https://api.github.com/repos/')
         DOM.Get(url,function(request){
           var data = JSON.parse(request.response);
          
           // if we got a reponse, try using it. 
           name.value =  data.name + " - " + data.description
           summary.value = data.description + " by " + data.owner.login
           
           // later use 
           // created_at -> original_published_at
           // updated_at -> original_updated_at
           // data.owner.name -> original_author
           
         },function(){
           console.log("failed to fetch github data")
         });
      
        return false;
      } 
      
      // We could also attempt to fetch the html page, and grab metadata from it 
      // author, pubdate, metadesc etc
      // would this be better done in a background way after story submission?
      
      
      // Else just use name from local url if we can
      name.value = urlToSentenceCase(original_url);
    
  
    }
  
    
    
  });

}

// Change a URL to a sentence for SetSubmitStoryName
function urlToSentenceCase(url) {
  if (url === undefined) {
    return ""
  }
  
  var parts, name
  url = url.replace(/\/$/, ""); // remove trailing /
  parts = url.split("/"); // now split on /
  name = parts[parts.length - 1]; // last part of string after last /
  name = name.replace(/[\?#].*$/, ""); //remove anything after ? or #
  name = name.replace(/^\d*-/, ""); // remove prefix numerals with dash (common on id based keys)
  name = name.replace(/\..*$/, ""); // remove .html etc extensions
  name = name.replace(/[_\-+]/g, " "); // remove all - or + or _ in string, replacing with space
  name = name.trim(); // remove whitespace trailing or leading
  name = name.toLowerCase(); // all lower
  name = name[0].toUpperCase() + name.substring(1); // Sentence case
  
  
  // Deal with some specific URLs
  if (url.match(/youtube|vimeo\.com/)) {
     name = "Video: "
  }
  if (url.match(/medium\.com/)) {
      // Eat the last word (UDID) on medium posts
      name = name.replace(/ [^ ]*$/, "");
  }

  
  
  return name
}