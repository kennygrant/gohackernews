// Package DOM provides functions to replace the use of jquery in 1.4KB of js
// See http://youmightnotneedjquery.com/ for more if required
var DOM = (function() {
    return {

        // Apply a function f on document ready
        Ready: function(f) {
            if (document.readyState != 'loading') {
                f();
            } else {
                document.addEventListener('DOMContentLoaded', f);
            }
        },

        // Return true if any elements match selector sel
        Exists: function(sel) {
            return (document.querySelector(sel) !== null);
        },

        // Return a NodeList of elements matching selector
        All: function(sel) {
            return document.querySelectorAll(sel);
        },

        // Return the first in the NodeList of elements matching selector - may return nil
        First: function(sel) {
            return DOM.All(sel)[0];
        },

        // Apply a function to elements of an array
        ForEach: function(array, f) {
            Array.prototype.forEach.call(array, f);
        },

        // Apply a function to elements matching selector, return true to break
        Each: function(sel, f) {
            var array = DOM.All(sel);
            for (i = 0; i < array.length; ++i) {
              f(array[i], i);
            }
        },

        // Attach event handlers to all matches for a selector 
        On: function(sel, name, f) {
            DOM.Each(sel, function(el, i) {
                el.addEventListener(name, f);
            });
        },
        
        // Return a NodeList of nearest elements matching selector, 
        // checking children, siblings or parents of el
        Nearest: function(el, sel) {

            // Start with this element, then walk up the tree till 
            // we find a child which matches selector or we run out of elements
            while (el !== undefined && el !== null) {
                var nearest = el.querySelectorAll(sel);
                if (nearest.length > 0) {
                    return nearest;
                }
                el = el.parentNode;
            }

            return []; // return empty array
        },

        // Attribute returns either an attribute value or an empty string (if null)
        Attribute: function(el, name) {
            if (el.getAttribute(name) === null) {
                return ''
            }
            return el.getAttribute(name)
        },

        // HasClass returns true if this element has this className
        HasClass: function(el, name) {
            var regexp = new RegExp("\\b" + name + "\\b", 'gi');
            return regexp.test(el.className);
        },

        // AddClass Adds the given className from el.className
        // s may be a string selector or an element
        AddClass: function(s, name) {
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    if (!DOM.HasClass(el, name)) {
                        el.className = el.className + ' ' + name;
                    }
                });
            } else {
                if (!DOM.HasClass(s, name)) {
                    s.className = s.className + ' ' + name;
                }
            }
        },

        // RemoveClass removes the given className from el.className
        // s may be a string selector or an element
        RemoveClass: function(s, name) {
            var regexp = new RegExp("\\b" + name + "\\b", 'gi');
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    el.className = el.className.replace(regexp, '')
                });
            } else {
                s.className = s.className.replace(regexp, '')
            }
        },

        // Format returns the format string with the indexed arguments substituted
        // Formats are of the form - "{0} {1}" which uses variables 0 and 1 respectively
        Format: function(format) {
            for (var i = 1; i < arguments.length; i++) {
                var regexp = new RegExp('\\{' + (i - 1) + '\\}', 'gi');
                format = format.replace(regexp, arguments[i]);
            }
            return format;
        },

        // Hide elements matching selector 
        // s may be a string selector or an element
        Hide: function(sel) {
            if (typeof sel === "string") {
                DOM.Each(sel, function(el, i) {
                    el.style.display = 'none';
                });
            } else {
                sel.style.display = 'none';
            }
        },

        // Show elements matching selector
        // s may be a string selector or an element
        Show: function(sel) {
            if (typeof sel === "string") {
                DOM.Each(sel, function(el, i) {
                    el.style.display = '';
                });
            } else {
                sel.style.display = '';
            }
        },

        // Hidden returns true if this element is hidden
        // s may be a string selector or an element
        Hidden: function(sel) {
            if (typeof sel === "string") {
                return (DOM.First(sel).style.display == 'none');
            } else {
                return sel.style.display == 'none';
            }

        },

        // Toggle the Shown or Hidden value of elements matching selector
        // s may be a string selector or an element
        ShowHide: function(sel) {
            if (typeof sel === "string") {
                DOM.Each(sel, function(el, i) {
                    if (el.style.display != 'none') {
                        el.style.display = 'none';
                    } else {
                        el.style.display = '';
                    }
                });
            } else {
                if (sel.style.display != 'none') {
                    sel.style.display = 'none';
                } else {
                    sel.style.display = '';
                }
            }
        },

        // Ajax - Get the data from url, call fs for success, fe for failures
        Get: function(url, fs, fe) {
            var request = new XMLHttpRequest();
            request.open('GET', url, true);
            request.onload = function() {
                if (request.status >= 200 && request.status < 400) {
                    fs(request);
                } else {
                    fe();
                }
            };
            request.onerror = fe;
            request.send();
        },

        // Ajax - post the data to url, call fs for success, fe for failures
        Post: function(url, data, fs, fe) {
            var request = new XMLHttpRequest();
            request.open('POST', url, true);
            request.onerror = fe;
            request.onload = function() {
                if (request.status >= 200 && request.status < 400) {
                    fs(request);
                } else {
                    fe(request);
                }
            };
            request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
            request.send(data);
        }

        

    };

}());