PD.PDTEXTController = Ember.ObjectController.extend({
 	welcome:'hi'
});

PD.PDLISTController = Ember.ArrayController.extend({
 	
});

PD.PDController = Ember.ObjectController.extend({
	actions: {
    createDoc: function () {
        // implement your action here
      var newdoc = this.store.createRecord('doc', {
        title: 'untitled',
        ctext: 'newdoc'
      });
      var that=this;
      newdoc.save().then();
      return false;
        //this.transitionToRoute('index');
    }
	}
});