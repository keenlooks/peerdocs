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
        toolbar: { container: '#toolbar-toolbar' }
      },
    theme: 'snow'
      });
      


      /*
      editor.on('selection-change', function(range) {
        console.log('selection-change', range)
      });
		*/

    var idInc=0;
    //,'doccgs':[{'location':33,'mod':'ssss'},{'location':33,'mod':'dddd'}]}
    var docdelt = this.get("store").createRecord('docdelt',{"docid":doc.get("id"),"cursor":0});
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
          var cg=that.get("store").createRecord('doccg',{'id':idInc,'location':(range.start-el.insert.length),'mod':el.insert});
          docdelt.get("doccgs").pushObject(cg);
          idInc++;
//docdelt.send('becomeDirty');
         //console.log(docdelt.serialize());
			   }
         else if(!freeze&&el.delete!=null && range!=null){
           console.log("delete: "+el.delete+" @ "+(range.start+el.delete)+" ** "+range.start);
           var deletestr="";
           for(var i=el.delete;i>0;i--){
              deletestr=deletestr+"\\b";
              //console.log(deletestr);
           }

        var cg2=that.get("store").createRecord('doccg',{'id':idInc,'location':(range.start+el.delete),'mod':deletestr});
         console.log(deletestr);
         deletestr="";
          docdelt.get("doccgs").pushObject(cg2);
          idInc++;
         }
			   else if(!freeze&&el.insert!=null){
				  console.log("insert :"+el.insert+" @ 0");
          var cg0=that.get("store").createRecord('doccg',{'id':idInc,'location':0,'mod':el.insert});
          docdelt.get("doccgs").pushObject(cg0);
          idInc++;
          //docdelt.send('becomeDirty');
           //console.log(docdelt.serialize());
			   }
		  });
        
    });

       refresher=setInterval(function(){

       // console.log(docdelt.serialize());
 
        console.log("00");
        var cursorp=0;
        if(editor.getSelection()!=null){
          cursorp=editor.getSelection().start;
        }
        docdelt.set("cursor",cursorp);
console.log("114");
        docdelt.save().then(function(docdelt){
          console.log("111");
          console.log(docdelt.get("doccgs").length);
          if(docdelt.get("doccgs").length!=0){

          docdelt.get("doccgs").clear();
          
          doc.reload().then(function(doc){
            console.log("222");
            editor.editor.disable();
            freeze=true;
            editor.setText(doc.get('ctext'));
            editor.setSelection(doc.get("cursor"),doc.get("cursor"));
            freeze=false;
            editor.editor.enable();
            console.log("333");

            console.log(doc.get('ctext')+"\n@"+doc.get("cursor"));
            
          });

          }

        });

        /*.then(function(docdelt){
          docdelt.clear();
          doc.reload();
        });*/

        },5000);


       editor.insertText(0, doc.get('ctext'), 'bold', true);
       editor.setSelection(doc.get("cursor"),doc.get("cursor"));
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


