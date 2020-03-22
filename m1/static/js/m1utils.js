// --------------------------------------------------------------------
// m1utils.js -- Various utilities for the m1 project
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------

function clearSelBox(sel) {
    var n = sel.length;
    var i;
    for (i = 0; i < n; i++) {
        sel.remove(0);
    }
}

function addselection(selbox, txt, lab, val) {
    var option = document.createElement("option");
    option.text = txt;
    option.label = lab;
    option.value = val;
    selbox.add(option);
}

function countMatchLen(s1, s2) {
  n = s1.length;
  if (n > s2.length) { n = s2.length;}
  for(var i = 0; i < n; i++) {
    c1 = s1.substr(i, 1).toLowerCase();
    c2 = s2.substr(i, 1).toLowerCase();
    if (c1 != c2) {
      return i;
    }
  }
  return n;
}

function setTextBox(fieldid, txt) {
    var d = document.getElementById(fieldid);
    d.text = txt;
    d.value = txt ;
}

function setSelectionBox(selid, txt) {
    //console.log("in setSelectioBox. Txt=", txt)
    var selbox = document.getElementById(selid);
    var i;
    // Look for an exact match, in either text or values.
    for(i = 0; i < selbox.options.length; i++) {
      var opt = selbox.options[i];
      if (opt.value == txt) {
          selbox.value = opt.value;
          //console.log("setting to (1):", selbox.value)
          return true;
      }
      if (opt.text == txt) {
          selbox.value = opt.value;
          //console.log("setting to (2):", selbox.value)
          return true;
      }
    }
    if (txt == "") {
        // This means that a correct default was not 
        // provided, so the current selection is valid.
        return true
    }
    // Return false if no valid selection was found.
    return false;
}

// fillUsers uses the JSON from "users" to fill a selection box 
function fillUsers(selbox_id, showall, add_blank=false) {
    dobj = document.getElementById("users");
    if (typeof dobj === "undefined") {
        console.log("fillUsers, users json unknown.")
        return;
    }
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillUsers, selbox id unknown: ", selbox_id)
        return;
    }
    var userlst = JSON.parse(dobj.innerHTML);
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for(i = 0; i < userlst.length; i++) {
        d = userlst[i];
        addselection(selbox, d.Name, d.Name, d.Name);
    }
}

