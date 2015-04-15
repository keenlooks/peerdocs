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
        ctext: 'newdoc'
      });
      var that=this;
      newdoc.save().then(function(saveddoc){that.transitionToRoute('PD.TEXT',saveddoc);}/*that.transitionToRoute('index')*/);
      return false;
      //+newdoc.get("id")
        //this.transitionToRoute('index');
    }
	}
});