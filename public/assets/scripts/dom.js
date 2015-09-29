// Package DOM provides functions to replace the use of jquery in 1.4KB of js - Ajax, Selectors, Event binding, ShowHide
// See http://youmightnotneedjquery.com/ for more if required
// Version 1.0.1
var DOM = (function() {
return {
  // Apply a function on document ready
  Ready:function(f) {
    if (document.readyState != 'loading'){
      f();
    } else {
      document.addEventListener('DOMContentLoaded', f);
    }
  },
  
  // Return true if any element match selector
  Exists:function(s) {
    return (document.querySelector(s) !== null);
  },

  // Return a NodeList of elements matching selector
  All:function(s) {
    return document.querySelectorAll(s);
  },
  
  // Return a NodeList of nearest elements matching selector, as children, siblings or parents of selector
  Nearest:function(el,s) {
    
    // Start with this element, then walk up the tree till we find a child which matches selector
    // Or we run out of elements
    while (el !== undefined) {
      var nearest = el.querySelectorAll(s);
      if (nearest.length > 0) {
        return nearest;
      }
      el = el.parentNode;
    }
    
    return undefined;
  },
  
  // Return the first in the NodeList of elements matching selector - may return undefined
  First:function(s) {
    return DOM.All(s)[0];
  },
  
  // Apply a function to elements matching selector
  Each:function(s,f) {
    var a = DOM.All(s);
    for (i = 0; i < a.length; ++i) {
      f(a[i],i);
    }
  },
  
  // Apply a function to elements of an array
  ForEach:function(a,f) {
    Array.prototype.forEach.call(a,f);
  },
  
  
  // Hide elements matching selector
  Hide:function(s) {
    DOM.Each(s,function(el,i){
      el.style.display = 'none';
    });
  },
  
  // Show elements matching selector
  Show:function(s) {
    DOM.Each(s,function(el,i){
      el.style.display = '';
    });
  },
  
  // Toggle the Shown or Hidden value of elements matching selector
  ShowHide:function(s) {
    DOM.Each(s,function(el,i){
      if (el.style.display != 'none') {
          el.Hide();
      } else {
         el.Show();
      }
    });
  },
  
  // Attach event handlers to all matches for a selector 
  On:function(s,b,f) {
    DOM.Each(s,function(el,i){
      el.addEventListener(b, f);
    });
  },
  
  // Format returns the format string with the indexed arguments substituted
  // Formats are of the form - "{0} {1}" which uses variables 0 and 1 respectively
  // TODO: We could at a later date perhaps accept named arguments?
  Format:function(f) {
    for (var i = 1; i < arguments.length; i++) {
      var regexp = new RegExp('\\{'+(i-1)+'\\}', 'gi');
      f = f.replace(regexp, arguments[i]);
    }
    return f;
  },
  
  // Ajax - Send the data d to url u, fs handles success, ff handles failures
  Post:function(u,d,fs,fe) {
    var request = new XMLHttpRequest();
    request.onreadystatechange = function(){
      if (request.readyState==4 && request.status==200) {
        fs(request);
      } else {
        fe(request);
      }
    }
    request.open('POST', u, true);
    request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
    request.send(d);
  },
  
  // Ajax - Get the data from url u, fs for success, ff for failures
  Get:function(u,fs,fe) {
    var request = new XMLHttpRequest();
    request.open('GET', u, true);
    request.onload = function() {
      if (request.status >= 200 && request.status < 400) {
        fs(request);
      } else {
        fe();
      }
    };
    request.onerror = fe;
    request.send();
    }
      
  };
    
}());