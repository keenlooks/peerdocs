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

    //,'doccgs':[{'location':33,'mod':'ssss'},{'location':33,'mod':'dddd'}]}
    var docdelt = this.get("store").createRecord('docdelt',{"docid":doc.get("id")});
    console.log(docdelt.serialize());
    docdelt.save();
    var freeze=false;
      var that=this;
      var cursorpos=0;
      editor.on('text-change', function(delta, source) {

        console.log("CHANGE FIRES "+freeze);
      	var range = editor.getSelection();

      	$.each(delta.ops,function(i,el){
          console.log(el+':'+range);
			   if(!freeze&&el.insert!=null && range!=null){
          cursorpos=range.start;
   				console.log("insert :"+el.insert+" @ "+(range.start-el.insert.length));
          var cg=that.get("store").createRecord('doccg',{'location':(range.start-el.insert.length),'mod':el.insert});
          docdelt.get("doccgs").pushObject(cg);
//docdelt.send('becomeDirty');
         //console.log(docdelt.serialize());
			   }
         else if(!freeze&&el.delete!=null && range!=null){
           console.log("delete: "+el.delete+" @ "+(range.start));
           var deletestr="";
           for(var i=el.delete;i>0;i--){
              deletestr=deletestr+"\\b";
              //console.log(deletestr);
           }

          var cg2=that.get("store").createRecord('doccg',{'location':range.start,'mod':deletestr});
         console.log(deletestr);
         deletestr="";
          docdelt.get("doccgs").pushObject(cg2);
         }
			   else if(!freeze&&el.insert!=null){
				  console.log("insert :"+el.insert+" @ 0");
          var cg0=that.get("store").createRecord('doccg',{'location':0,'mod':el.insert});
          docdelt.get("doccgs").pushObject(cg0);
          //docdelt.send('becomeDirty');
           //console.log(docdelt.serialize());
			   }
		  });
        
    });

       refresher=setInterval(function(){

       // console.log(docdelt.serialize());
 
        //console.log("hi");
        //docdelt.set("docid",doc.get("id"));
        docdelt.save().then(function(docdelt){
          docdelt.get("doccgs").clear();
          //console.log(docdelt.serialize());
          doc.reload().then(function(doc){

            freeze=true;
            editor.setText(doc.get('ctext'));
            freeze=false;

            console.log(doc.get('ctext'));
            editor.setSelection(cursorpos,cursorpos);
          });
          

        });

        /*.then(function(docdelt){
          docdelt.clear();
          doc.reload();
        });*/

        },5000);


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


