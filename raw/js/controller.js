PD.PDIndexController = Ember.ObjectController.extend({
actions: {
  shownav:function(){
    $("#opennav").click(function(){
      $("#services").show();
      return false;
    })
  }
  }
});


PD.PDTEXTController = Ember.ObjectController.extend({
 	welcome:'hi'
});

PD.PDLISTController = Ember.ArrayController.extend({
 	
});

PD.PDController = Ember.ObjectController.extend({
	actions: {

  shownav:function(){
    $("#opennav").click(function(){
      $("#services").show();
      return false;
    })
  },

  createDoc: function () {
        // implement your action here
      var newdoc = this.store.createRecord('doc', {
        title: 'untitled',
        ctext: 'newdoc',
        cursor: 0
      });
      var that=this;
      newdoc.save().then(function(saveddoc){that.transitionToRoute('PD.TEXT',saveddoc);}/*that.transitionToRoute('index')*/);
      return false;
      //+newdoc.get("id")
        //this.transitionToRoute('index');
    },
  invite:function(){
    var invitation = this.store.createRecord('invitation', {
        address: $("#iip").val(),
        dockey: $("#ikey").val(),
        name: $("#iname").val(),
        docid:parseInt($("#idoc").val(),10),
        type:"invite"
      });
    console.log("1");
    invitation.save();
 console.log("2");
    return false;

  },

  join:function(){
    var invitation = this.store.createRecord('invitation', {
        //address: $("#iip").val(),
        //dockey: $("#ikey").val(),
       // name: $("#iname").val(),
        docid:parseInt($("#idoc").val(),10),
        type:"join"
      });
    
    invitation.save();

    return false;

  }
	}
});