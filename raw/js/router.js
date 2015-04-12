PD.Router.map(function() {
  this.resource('PD', { path: '/' },function(){
  	this.route('TEXT',{path:'TEXT/:TEXT_id'});
  	this.route('LIST');
  	this.route('SETTING');
  	this.route('INVITE');
  });

});

var refresher;

PD.PDTEXTRoute = Ember.Route.extend({

model: function(params) {
    return this.store.find('doc',params.TEXT_id);
},


setupController: function(controller,doc){
	   controller.set('model', doc);

       Ember.run.schedule('afterRender', this, function(){
       

      var editor = new Quill('#editor-container', {
        modules: {
          'toolbar': { container: '#formatting-container' },
          'link-tooltip': true,
          'image-tooltip': true
        }
      });
      /*
      editor.on('selection-change', function(range) {
        console.log('selection-change', range)
      });
		*/
    var docdelt = this.store.createRecord('docdelt');
    //console.log(docdelta);
    docdelt.save();

      var that=this;
      editor.on('text-change', function(delta, source) {
      	var range = editor.getSelection();

      	$.each(delta.ops,function(i,el){
			   if(el.insert!=null && range!=null){
   				console.log(el.insert+" @ "+(range.start-el.insert.length));
          var cg=that.store.createRecord('doccg',{'location':(range.start-el.insert.length),'mod':el.insert});
          docdelt.get("doccgs").addObject(cg);
          //console.log(docdelta);
			   }
			   else if(el.insert!=null){
				  console.log(el.insert+" @ 0");
          var cg0=that.store.createRecord('doccg',{'location':0,'mod':el.insert});
          docdelt.get("doccgs").addObject(cg0);
          //console.log(docdelta);
			   }
		  });
        
    });

       refresher=setInterval(function(){
        //console.log("hi");
        docdelt.save().then(function(docdelt){
          docdelt.reload();
        });

        },3000);


       editor.insertText(0, doc.get('ctext'), 'bold', true);
    });


  },
  renderTemplate: function(controller) {
    this.render('pd/TEXT', {controller:controller});
  },

  deactivate:function() {
    clearInterval(refresher);
  }
});

PD.PDLISTRoute = Ember.Route.extend({
model: function() {

    return this.store.find('docmeta');
},

setupController: function(controller,docmeta){
	   controller.set('model', docmeta);
},

renderTemplate: function(controller) {
    this.render('pd/LIST', {controller: controller});
  }
});

PD.PDSETTINGRoute = Ember.Route.extend({
  renderTemplate: function(controller) {
    this.render('pd/NI', {controller: controller});
  }
});

PD.PDINVITERoute = Ember.Route.extend({
  renderTemplate: function(controller) {
    this.render('pd/NI', {controller: controller});
  }
});


