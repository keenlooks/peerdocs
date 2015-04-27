/* jshint ignore:start */

/* jshint ignore:end */

define('peerdoc-embercli/adapters/application', ['exports', 'ember-data'], function (exports, DS) {

  'use strict';

  exports['default'] = DS['default'].RESTAdapter.extend({
    namespace: 'api',
    host: 'http://localhost:8080' });

});
define('peerdoc-embercli/app', ['exports', 'ember', 'ember/resolver', 'ember/load-initializers', 'peerdoc-embercli/config/environment'], function (exports, Ember, Resolver, loadInitializers, config) {

  'use strict';

  var App;

  Ember['default'].MODEL_FACTORY_INJECTIONS = true;

  App = Ember['default'].Application.extend({
    modulePrefix: config['default'].modulePrefix,
    podModulePrefix: config['default'].podModulePrefix,
    Resolver: Resolver['default']
  });

  loadInitializers['default'](App, config['default'].modulePrefix);

  exports['default'] = App;

});
define('peerdoc-embercli/controllers/docmeta', ['exports', 'ember'], function (exports, Ember) {

  'use strict';

  exports['default'] = Ember['default'].ObjectController.extend({

    isModified: (function () {
      var model = this.get('model');

      //return model.get('lastmod')=="true";
      return true;
    }).property('model.lastmod'),

    isPending: (function () {
      var model = this.get('model');

      //return model.get('lastmod')=="pending";  
      return true;
    }).property('model.lastmod') });

});
define('peerdoc-embercli/controllers/pd', ['exports', 'ember'], function (exports, Ember) {

    'use strict';

    exports['default'] = Ember['default'].ArrayController.extend({
        itemController: 'docmeta',
        sortProperties: ['id'],
        sortAscending: false,
        actions: {
            createDoc: function createDoc() {
                // implement your action here
                var model = this.get('model');
                var title = this.get('newTitle');
                if (!title.trim()) {
                    return;
                }

                var newdoc = this.store.createRecord('doc', {
                    title: title,
                    ctext: 'newdoc',
                    cursor: 0
                });
                var that = this;
                newdoc.save().then(function (saveddoc) {
                    model.update().then(function () {

                        that.transitionToRoute('pd.doc', saveddoc);
                    });
                } /*that.transitionToRoute('index')*/);
                $('#new-todo').val('');

                return false;
                //+newdoc.get("id")
                //this.transitionToRoute('index');
            },

            invite: function invite(id, istr) {
                if (!istr) {
                    console.log('NOT CORRECT INPUT FORMAT');
                    return false;
                }
                var strs = istr.split('@');
                if (strs.length != 2) {
                    console.log('NOT CORRECT INPUT FORMAT');
                    return false;
                }
                var invitation = this.store.createRecord('invitation', {
                    address: strs[0],
                    dockey: 'xxxxx',
                    name: strs[1],
                    docid: id,
                    type: 'invite'
                });
                //console.log("1");
                invitation.save();
                //console.log("2");
                return false;
            },

            join: function join(id) {

                var joinvar = this.store.createRecord('invitation', {
                    dockey: 'xxxxx',
                    docid: id,
                    type: 'join'
                });
                //console.log("1");
                joinvar.save();
                //console.log("2");
                return false;
            } }
    });

});
define('peerdoc-embercli/controllers/pd/doc', ['exports', 'ember'], function (exports, Ember) {

   'use strict';

   exports['default'] = Ember['default'].Controller.extend({
      needs: "pd" });

});
define('peerdoc-embercli/controllers/pd/index', ['exports', 'ember'], function (exports, Ember) {

	'use strict';

	exports['default'] = Ember['default'].Controller.extend({});

});
define('peerdoc-embercli/helpers/if-cond', ['exports', 'ember'], function (exports, Ember) {

  'use strict';

  exports.ifCond = ifCond;

  function ifCond(v1, v2, options) {

    if (v1 === v2) {
      return options.fn(this);
    }
    return options.inverse(this);
  }

  exports['default'] = Ember['default'].HTMLBars.makeBoundHelper(ifCond);

});
define('peerdoc-embercli/initializers/app-version', ['exports', 'peerdoc-embercli/config/environment', 'ember'], function (exports, config, Ember) {

  'use strict';

  var classify = Ember['default'].String.classify;
  var registered = false;

  exports['default'] = {
    name: 'App Version',
    initialize: function initialize(container, application) {
      if (!registered) {
        var appName = classify(application.toString());
        Ember['default'].libraries.register(appName, config['default'].APP.version);
        registered = true;
      }
    }
  };

});
define('peerdoc-embercli/initializers/export-application-global', ['exports', 'ember', 'peerdoc-embercli/config/environment'], function (exports, Ember, config) {

  'use strict';

  exports.initialize = initialize;

  function initialize(container, application) {
    var classifiedName = Ember['default'].String.classify(config['default'].modulePrefix);

    if (config['default'].exportApplicationGlobal && !window[classifiedName]) {
      window[classifiedName] = application;
    }
  }

  ;

  exports['default'] = {
    name: 'export-application-global',

    initialize: initialize
  };

});
define('peerdoc-embercli/models/doc', ['exports', 'ember-data'], function (exports, DS) {

  'use strict';

  exports['default'] = DS['default'].Model.extend({
    title: DS['default'].attr('string'),
    ctext: DS['default'].attr('string'),
    cursor: DS['default'].attr('number')
  });

});
define('peerdoc-embercli/models/doccg', ['exports', 'ember-data'], function (exports, DS) {

    'use strict';

    exports['default'] = DS['default'].Model.extend({
        location: DS['default'].attr('number'),
        mod: DS['default'].attr('string')
    });

});
define('peerdoc-embercli/models/docdelt', ['exports', 'ember-data'], function (exports, DS) {

  'use strict';

  exports['default'] = DS['default'].Model.extend({
    docid: DS['default'].attr('number'),
    cursor: DS['default'].attr('number'),
    doccgs: DS['default'].hasMany('doccg', { embedded: 'always' }) });

});
define('peerdoc-embercli/models/docmeta', ['exports', 'ember-data'], function (exports, DS) {

   'use strict';

   exports['default'] = DS['default'].Model.extend({
      title: DS['default'].attr('string'),
      lastmod: DS['default'].attr('string')
   });

});
define('peerdoc-embercli/models/invitation', ['exports', 'ember-data'], function (exports, DS) {

  'use strict';

  exports['default'] = DS['default'].Model.extend({
    docid: DS['default'].attr('number'),
    address: DS['default'].attr('string'),
    type: DS['default'].attr('string'),
    name: DS['default'].attr('string'),
    dockey: DS['default'].attr('string')
  });

});
define('peerdoc-embercli/models/pd', ['exports', 'ember-data'], function (exports, DS) {

	'use strict';

	exports['default'] = DS['default'].Model.extend({});

});
define('peerdoc-embercli/router', ['exports', 'ember', 'peerdoc-embercli/config/environment'], function (exports, Ember, config) {

  'use strict';

  var Router = Ember['default'].Router.extend({
    location: config['default'].locationType
  });

  exports['default'] = Router.map(function () {
    //this.resource('doc',{path:'docs/:doc_id'});

    this.resource('pd', { path: '/' }, function () {
      this.route('doc', { path: 'docs/:doc_id' });
      //this.route('LIST');
      //this.route('SETTING');
      //this.route('INVITE');
    });
  });

});
define('peerdoc-embercli/routes/pd', ['exports', 'ember'], function (exports, Ember) {

  'use strict';

  exports['default'] = Ember['default'].Route.extend({

    refresher: null,

    model: function model() {
      return this.store.find('docmeta');
    },

    clearAllTimer: function clearAllTimer() {
      for (var i = 0; i < 5000; i++) {
        if (this.refresher == null || i != this.refresher) {
          clearInterval(i);
        }
      }
    },

    registerInvite: function registerInvite() {

      var that = this;

      $.each($('.glyphicon-plus').toArray(), function (i, o) {
        var ev = $._data(o, 'events');
        if (ev == null) {
          $('.glyphicon-plus').eq(i).click(function () {

            $(this).parent().next().slideToggle();
          });

          $('.glyphicon-plus').eq(i).hover(function () {
            $(this).removeClass('glyphicon-plus').addClass('glyphicon-remove');
          }, function () {
            $(this).removeClass('glyphicon-remove').addClass('glyphicon-plus');
          });

          $('.glyphicon-plus').eq(i).parent().next().find('.invite-input').focus(function () {
            that.clearAllTimer();
          });
        }
      });
    },

    updateMeta: function updateMeta(docmeta) {
      var that = this;
      return setInterval(function () {
        console.log('updatemeta');
        docmeta.update();
        that.registerInvite();
      }, 1000);
    },

    setupController: function setupController(controller, docmeta) {
      var that = this;
      controller.set('model', docmeta);

      controller.set('test', 'test');

      Ember['default'].run.schedule('afterRender', this, function () {

        that.refresher = that.updateMeta(docmeta);
        controller.set('refresher', that.refresher);
        console.log('MR :' + that.refresher);

        $('.glyphicon-plus').click(function () {

          $(this).parent().next().slideToggle();
        });

        $('.glyphicon-plus').hover(function () {
          $(this).removeClass('glyphicon-plus').addClass('glyphicon-remove');
        }, function () {
          $(this).removeClass('glyphicon-remove').addClass('glyphicon-plus');
        });

        $('#new-todo,.invite-input').focus(function () {
          that.clearAllTimer();
        });

        $('#findm').click(function () {
          $('#findm').hide();
          $('#nav-icon').show();
          $('.sidebar').animate({ width: 'toggle' }, 350);
        });

        $('#nav-icon').click(function () {

          $('.sidebar').animate({ width: 'toggle' }, 350);

          if ($('#nav-icon').attr('class') == 'nav-icon-left') {
            $('#nav-icon').attr('class', 'nav-icon-right');
          } else {
            $('#nav-icon').attr('class', 'nav-icon-left');
          }
        });
      }); //afterRender
    }, //set up controller

    renderTemplate: function renderTemplate(controller) {

      this.render('pd');
      //this.render("welcome",{output:"main",into:"pd"});
    }

  });

});
define('peerdoc-embercli/routes/pd/doc', ['exports', 'ember'], function (exports, Ember) {

  'use strict';

  exports['default'] = Ember['default'].Route.extend({

    model: function model(params) {
      return this.store.find('doc', params.doc_id);
    },

    setupController: function setupController(controller, doc) {
      controller.set('model', doc);
      controller.set('title', doc.get('title'));
      console.log(doc.get('title'));
      var refresher;
      console.log('MR :' + controller.get('controllers.pd.refresher'));
      for (var i = 0; i < 5000; i++) {
        if (controller.get('controllers.pd.refresher') == null || i != controller.get('controllers.pd.refresher')) {
          clearInterval(i);
        }
      }

      Ember['default'].run.schedule('afterRender', this, function () {

        var editor = new Quill('#editor-container');

        controller.set('refresher', refresher);
        /*
        editor.on('selection-change', function(range) {
          console.log('selection-change', range)
        });
        */

        var idInc = Math.floor(Math.random() * 10000 + 1);
        //,'doccgs':[{'location':33,'mod':'ssss'},{'location':33,'mod':'dddd'}]}
        var docdelt = this.get('store').createRecord('docdelt', { docid: doc.get('id'), cursor: 0 });
        console.log(docdelt.serialize());
        docdelt.save();
        var freeze = false;
        var that = this;
        var cursorpos = 0;
        editor.on('text-change', function (delta, source) {

          console.log('CHANGE FIRES ' + freeze);
          var range = editor.getSelection();

          $.each(delta.ops, function (i, el) {
            console.log(el + ':' + range);
            if (!freeze && el.insert != null && range != null) {
              cursorpos = range.start;
              console.log('insert :' + el.insert + ' @ ' + (range.start - el.insert.length));
              var cg = that.get('store').createRecord('doccg', { id: idInc, location: range.start - el.insert.length, mod: el.insert });
              docdelt.get('doccgs').pushObject(cg);
              idInc++;
              //docdelt.send('becomeDirty');
              //console.log(docdelt.serialize());
            } else if (!freeze && el['delete'] != null && range != null) {
              console.log('delete: ' + el['delete'] + ' @ ' + (range.start + el['delete']) + ' ** ' + range.start);
              var deletestr = '';
              for (var i = el['delete']; i > 0; i--) {
                deletestr = deletestr + '\\b';
                //console.log(deletestr);
              }

              var cg2 = that.get('store').createRecord('doccg', { id: idInc, location: range.start + el['delete'], mod: deletestr });
              console.log(deletestr);
              deletestr = '';
              docdelt.get('doccgs').pushObject(cg2);
              idInc++;
            } else if (!freeze && el.insert != null) {
              console.log('insert :' + el.insert + ' @ 0');
              var cg0 = that.get('store').createRecord('doccg', { id: idInc, location: 0, mod: el.insert });
              docdelt.get('doccgs').pushObject(cg0);
              idInc++;
              //docdelt.send('becomeDirty');
              //console.log(docdelt.serialize());
            }
          });
        });

        refresher = setInterval(function () {

          var cursorp = 0;
          if (editor.getSelection() != null) {
            cursorp = editor.getSelection().start;
          }
          docdelt.set('cursor', cursorp);

          docdelt.save().then(function (docdelt) {

            //console.log(docdelt.get("doccgs").size());
            console.log('PUT :' + docdelt.get('cursor') + ':' + docdelt.get('doccgs').length);
            if (docdelt.get('doccgs') != null) {

              docdelt.get('doccgs').clear();

              doc.reload().then(function (doc) {
                //console.log("222");
                if (doc.get('title') != 'None') {
                  editor.editor.disable();
                  freeze = true;
                  editor.setText(doc.get('ctext'));
                  if (doc.get('cursor') >= 0) {
                    editor.setSelection(doc.get('cursor'), doc.get('cursor'));
                  }
                  freeze = false;
                  editor.editor.enable();
                  //console.log("333");

                  console.log(doc.get('ctext') + '\n@' + doc.get('cursor'));
                }
              });
            }
          });

          /*.then(function(docdelt){
            docdelt.clear();
            doc.reload();
          });*/
        }, 300);

        editor.setText(doc.get('ctext'));
        editor.setSelection(doc.get('cursor'), doc.get('cursor'));
        docdelt.get('doccgs').clear();
        console.log('Timer' + refresher);
      });
    },
    renderTemplate: function renderTemplate(controller) {
      //this.render("doc");
      this.render('doc', { outet: 'main', into: 'pd' });
    },

    deactivate: function deactivate() {
      clearInterval(refresher);
    }
  });

});
define('peerdoc-embercli/serializers/docdelt', ['exports', 'ember-data'], function (exports, DS) {

  'use strict';

  exports['default'] = DS['default'].RESTSerializer.extend(DS['default'].EmbeddedRecordsMixin, {
    attrs: {
      doccgs: { embedded: 'always' }
    }
  });

});
define('peerdoc-embercli/templates/application', ['exports'], function (exports) {

  'use strict';

  exports['default'] = Ember.HTMLBars.template((function() {
    return {
      isHTMLBars: true,
      revision: "Ember@1.11.3",
      blockParams: 0,
      cachedFragment: null,
      hasRendered: false,
      build: function build(dom) {
        var el0 = dom.createDocumentFragment();
        var el1 = dom.createTextNode("\n    ");
        dom.appendChild(el0, el1);
        var el1 = dom.createElement("div");
        dom.setAttribute(el1,"class","site-wrapper");
        var el2 = dom.createTextNode("\n        ");
        dom.appendChild(el1, el2);
        var el2 = dom.createComment("");
        dom.appendChild(el1, el2);
        var el2 = dom.createTextNode("\n    ");
        dom.appendChild(el1, el2);
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n");
        dom.appendChild(el0, el1);
        return el0;
      },
      render: function render(context, env, contextualElement) {
        var dom = env.dom;
        var hooks = env.hooks, content = hooks.content;
        dom.detectNamespace(contextualElement);
        var fragment;
        if (env.useFragmentCache && dom.canClone) {
          if (this.cachedFragment === null) {
            fragment = this.build(dom);
            if (this.hasRendered) {
              this.cachedFragment = fragment;
            } else {
              this.hasRendered = true;
            }
          }
          if (this.cachedFragment) {
            fragment = dom.cloneNode(this.cachedFragment, true);
          }
        } else {
          fragment = this.build(dom);
        }
        var morph0 = dom.createMorphAt(dom.childAt(fragment, [1]),1,1);
        content(env, morph0, context, "outlet");
        return fragment;
      }
    };
  }()));

});
define('peerdoc-embercli/templates/doc', ['exports'], function (exports) {

  'use strict';

  exports['default'] = Ember.HTMLBars.template((function() {
    return {
      isHTMLBars: true,
      revision: "Ember@1.11.3",
      blockParams: 0,
      cachedFragment: null,
      hasRendered: false,
      build: function build(dom) {
        var el0 = dom.createDocumentFragment();
        var el1 = dom.createElement("span");
        dom.setAttribute(el1,"id","sticker");
        var el2 = dom.createComment("");
        dom.appendChild(el1, el2);
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n      ");
        dom.appendChild(el0, el1);
        var el1 = dom.createElement("div");
        dom.setAttribute(el1,"id","editor-wrapper");
        var el2 = dom.createTextNode("\n        ");
        dom.appendChild(el1, el2);
        var el2 = dom.createElement("div");
        dom.setAttribute(el2,"id","editor-container");
        dom.appendChild(el1, el2);
        var el2 = dom.createTextNode("\n      \n  ");
        dom.appendChild(el1, el2);
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n");
        dom.appendChild(el0, el1);
        var el1 = dom.createComment("");
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n");
        dom.appendChild(el0, el1);
        return el0;
      },
      render: function render(context, env, contextualElement) {
        var dom = env.dom;
        var hooks = env.hooks, content = hooks.content;
        dom.detectNamespace(contextualElement);
        var fragment;
        if (env.useFragmentCache && dom.canClone) {
          if (this.cachedFragment === null) {
            fragment = this.build(dom);
            if (this.hasRendered) {
              this.cachedFragment = fragment;
            } else {
              this.hasRendered = true;
            }
          }
          if (this.cachedFragment) {
            fragment = dom.cloneNode(this.cachedFragment, true);
          }
        } else {
          fragment = this.build(dom);
        }
        var morph0 = dom.createMorphAt(dom.childAt(fragment, [0]),0,0);
        var morph1 = dom.createMorphAt(fragment,4,4,contextualElement);
        content(env, morph0, context, "title");
        content(env, morph1, context, "outlet");
        return fragment;
      }
    };
  }()));

});
define('peerdoc-embercli/templates/pd', ['exports'], function (exports) {

  'use strict';

  exports['default'] = Ember.HTMLBars.template((function() {
    var child0 = (function() {
      var child0 = (function() {
        return {
          isHTMLBars: true,
          revision: "Ember@1.11.3",
          blockParams: 0,
          cachedFragment: null,
          hasRendered: false,
          build: function build(dom) {
            var el0 = dom.createDocumentFragment();
            var el1 = dom.createTextNode("            ");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("span");
            dom.setAttribute(el1,"class","doc-state banner_m");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n");
            dom.appendChild(el0, el1);
            return el0;
          },
          render: function render(context, env, contextualElement) {
            var dom = env.dom;
            dom.detectNamespace(contextualElement);
            var fragment;
            if (env.useFragmentCache && dom.canClone) {
              if (this.cachedFragment === null) {
                fragment = this.build(dom);
                if (this.hasRendered) {
                  this.cachedFragment = fragment;
                } else {
                  this.hasRendered = true;
                }
              }
              if (this.cachedFragment) {
                fragment = dom.cloneNode(this.cachedFragment, true);
              }
            } else {
              fragment = this.build(dom);
            }
            return fragment;
          }
        };
      }());
      var child1 = (function() {
        return {
          isHTMLBars: true,
          revision: "Ember@1.11.3",
          blockParams: 0,
          cachedFragment: null,
          hasRendered: false,
          build: function build(dom) {
            var el0 = dom.createDocumentFragment();
            var el1 = dom.createTextNode("     \n");
            dom.appendChild(el0, el1);
            return el0;
          },
          render: function render(context, env, contextualElement) {
            var dom = env.dom;
            dom.detectNamespace(contextualElement);
            var fragment;
            if (env.useFragmentCache && dom.canClone) {
              if (this.cachedFragment === null) {
                fragment = this.build(dom);
                if (this.hasRendered) {
                  this.cachedFragment = fragment;
                } else {
                  this.hasRendered = true;
                }
              }
              if (this.cachedFragment) {
                fragment = dom.cloneNode(this.cachedFragment, true);
              }
            } else {
              fragment = this.build(dom);
            }
            return fragment;
          }
        };
      }());
      var child2 = (function() {
        return {
          isHTMLBars: true,
          revision: "Ember@1.11.3",
          blockParams: 0,
          cachedFragment: null,
          hasRendered: false,
          build: function build(dom) {
            var el0 = dom.createDocumentFragment();
            var el1 = dom.createTextNode("   ");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("span");
            dom.setAttribute(el1,"class","doc-state banner_s");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n    ");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("a");
            dom.setAttribute(el1,"href","#");
            var el2 = dom.createComment("");
            dom.appendChild(el1, el2);
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n    ");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("span");
            dom.setAttribute(el1,"class","glyphicon white glyphicon-user right");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n");
            dom.appendChild(el0, el1);
            return el0;
          },
          render: function render(context, env, contextualElement) {
            var dom = env.dom;
            var hooks = env.hooks, content = hooks.content, get = hooks.get, element = hooks.element;
            dom.detectNamespace(contextualElement);
            var fragment;
            if (env.useFragmentCache && dom.canClone) {
              if (this.cachedFragment === null) {
                fragment = this.build(dom);
                if (this.hasRendered) {
                  this.cachedFragment = fragment;
                } else {
                  this.hasRendered = true;
                }
              }
              if (this.cachedFragment) {
                fragment = dom.cloneNode(this.cachedFragment, true);
              }
            } else {
              fragment = this.build(dom);
            }
            var element3 = dom.childAt(fragment, [5]);
            var morph0 = dom.createMorphAt(dom.childAt(fragment, [3]),0,0);
            content(env, morph0, context, "docmeta.title");
            element(env, element3, context, "action", ["join", get(env, context, "docmeta.id")], {});
            return fragment;
          }
        };
      }());
      var child3 = (function() {
        var child0 = (function() {
          return {
            isHTMLBars: true,
            revision: "Ember@1.11.3",
            blockParams: 0,
            cachedFragment: null,
            hasRendered: false,
            build: function build(dom) {
              var el0 = dom.createDocumentFragment();
              var el1 = dom.createComment("");
              dom.appendChild(el0, el1);
              return el0;
            },
            render: function render(context, env, contextualElement) {
              var dom = env.dom;
              var hooks = env.hooks, content = hooks.content;
              dom.detectNamespace(contextualElement);
              var fragment;
              if (env.useFragmentCache && dom.canClone) {
                if (this.cachedFragment === null) {
                  fragment = this.build(dom);
                  if (this.hasRendered) {
                    this.cachedFragment = fragment;
                  } else {
                    this.hasRendered = true;
                  }
                }
                if (this.cachedFragment) {
                  fragment = dom.cloneNode(this.cachedFragment, true);
                }
              } else {
                fragment = this.build(dom);
              }
              var morph0 = dom.createMorphAt(fragment,0,0,contextualElement);
              dom.insertBoundary(fragment, null);
              dom.insertBoundary(fragment, 0);
              content(env, morph0, context, "docmeta.title");
              return fragment;
            }
          };
        }());
        return {
          isHTMLBars: true,
          revision: "Ember@1.11.3",
          blockParams: 0,
          cachedFragment: null,
          hasRendered: false,
          build: function build(dom) {
            var el0 = dom.createDocumentFragment();
            var el1 = dom.createTextNode("    ");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("span");
            dom.setAttribute(el1,"class","doc-state hidden");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n    ");
            dom.appendChild(el0, el1);
            var el1 = dom.createComment("");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("span");
            dom.setAttribute(el1,"class","glyphicon white glyphicon-plus right");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n");
            dom.appendChild(el0, el1);
            return el0;
          },
          render: function render(context, env, contextualElement) {
            var dom = env.dom;
            var hooks = env.hooks, get = hooks.get, block = hooks.block;
            dom.detectNamespace(contextualElement);
            var fragment;
            if (env.useFragmentCache && dom.canClone) {
              if (this.cachedFragment === null) {
                fragment = this.build(dom);
                if (this.hasRendered) {
                  this.cachedFragment = fragment;
                } else {
                  this.hasRendered = true;
                }
              }
              if (this.cachedFragment) {
                fragment = dom.cloneNode(this.cachedFragment, true);
              }
            } else {
              fragment = this.build(dom);
            }
            var morph0 = dom.createMorphAt(fragment,3,3,contextualElement);
            block(env, morph0, context, "link-to", ["pd.doc", get(env, context, "docmeta.id")], {}, child0, null);
            return fragment;
          }
        };
      }());
      var child4 = (function() {
        return {
          isHTMLBars: true,
          revision: "Ember@1.11.3",
          blockParams: 0,
          cachedFragment: null,
          hasRendered: false,
          build: function build(dom) {
            var el0 = dom.createDocumentFragment();
            var el1 = dom.createTextNode(" \n");
            dom.appendChild(el0, el1);
            return el0;
          },
          render: function render(context, env, contextualElement) {
            var dom = env.dom;
            dom.detectNamespace(contextualElement);
            var fragment;
            if (env.useFragmentCache && dom.canClone) {
              if (this.cachedFragment === null) {
                fragment = this.build(dom);
                if (this.hasRendered) {
                  this.cachedFragment = fragment;
                } else {
                  this.hasRendered = true;
                }
              }
              if (this.cachedFragment) {
                fragment = dom.cloneNode(this.cachedFragment, true);
              }
            } else {
              fragment = this.build(dom);
            }
            return fragment;
          }
        };
      }());
      var child5 = (function() {
        return {
          isHTMLBars: true,
          revision: "Ember@1.11.3",
          blockParams: 0,
          cachedFragment: null,
          hasRendered: false,
          build: function build(dom) {
            var el0 = dom.createDocumentFragment();
            var el1 = dom.createTextNode("\n");
            dom.appendChild(el0, el1);
            var el1 = dom.createElement("div");
            dom.setAttribute(el1,"class","invite");
            var el2 = dom.createTextNode("\n ");
            dom.appendChild(el1, el2);
            var el2 = dom.createElement("div");
            dom.setAttribute(el2,"class","input-group");
            var el3 = dom.createTextNode("\n\n      ");
            dom.appendChild(el2, el3);
            var el3 = dom.createComment("");
            dom.appendChild(el2, el3);
            var el3 = dom.createTextNode("\n      ");
            dom.appendChild(el2, el3);
            var el3 = dom.createElement("span");
            dom.setAttribute(el3,"class","input-group-btn");
            var el4 = dom.createTextNode("\n        ");
            dom.appendChild(el3, el4);
            var el4 = dom.createElement("button");
            dom.setAttribute(el4,"class","btn btn-default");
            dom.setAttribute(el4,"type","button");
            var el5 = dom.createTextNode("Go!");
            dom.appendChild(el4, el5);
            dom.appendChild(el3, el4);
            var el4 = dom.createTextNode("\n      ");
            dom.appendChild(el3, el4);
            dom.appendChild(el2, el3);
            var el3 = dom.createTextNode("\n    ");
            dom.appendChild(el2, el3);
            dom.appendChild(el1, el2);
            var el2 = dom.createComment(" /input-group ");
            dom.appendChild(el1, el2);
            var el2 = dom.createTextNode("\n\n");
            dom.appendChild(el1, el2);
            var el2 = dom.createElement("span");
            dom.setAttribute(el2,"style","display:none");
            var el3 = dom.createComment("");
            dom.appendChild(el2, el3);
            dom.appendChild(el1, el2);
            var el2 = dom.createTextNode("\n\n");
            dom.appendChild(el1, el2);
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode(" ");
            dom.appendChild(el0, el1);
            var el1 = dom.createComment("invite");
            dom.appendChild(el0, el1);
            var el1 = dom.createTextNode("\n\n");
            dom.appendChild(el0, el1);
            return el0;
          },
          render: function render(context, env, contextualElement) {
            var dom = env.dom;
            var hooks = env.hooks, get = hooks.get, inline = hooks.inline, element = hooks.element, content = hooks.content;
            dom.detectNamespace(contextualElement);
            var fragment;
            if (env.useFragmentCache && dom.canClone) {
              if (this.cachedFragment === null) {
                fragment = this.build(dom);
                if (this.hasRendered) {
                  this.cachedFragment = fragment;
                } else {
                  this.hasRendered = true;
                }
              }
              if (this.cachedFragment) {
                fragment = dom.cloneNode(this.cachedFragment, true);
              }
            } else {
              fragment = this.build(dom);
            }
            var element0 = dom.childAt(fragment, [1]);
            var element1 = dom.childAt(element0, [1]);
            var element2 = dom.childAt(element1, [3, 1]);
            var morph0 = dom.createMorphAt(element1,1,1);
            var morph1 = dom.createMorphAt(dom.childAt(element0, [4]),0,0);
            inline(env, morph0, context, "input", [], {"type": "text", "class": "form-control invite-input", "placeholder": "name@ip_addr", "value": get(env, context, "istr")});
            element(env, element2, context, "action", ["invite", get(env, context, "docmeta.id"), get(env, context, "istr")], {});
            content(env, morph1, context, "this.id");
            return fragment;
          }
        };
      }());
      return {
        isHTMLBars: true,
        revision: "Ember@1.11.3",
        blockParams: 0,
        cachedFragment: null,
        hasRendered: false,
        build: function build(dom) {
          var el0 = dom.createDocumentFragment();
          var el1 = dom.createTextNode(" \n          ");
          dom.appendChild(el0, el1);
          var el1 = dom.createElement("li");
          var el2 = dom.createTextNode("\n");
          dom.appendChild(el1, el2);
          var el2 = dom.createComment("");
          dom.appendChild(el1, el2);
          var el2 = dom.createTextNode("\n		    \n");
          dom.appendChild(el1, el2);
          var el2 = dom.createComment("");
          dom.appendChild(el1, el2);
          dom.appendChild(el0, el1);
          var el1 = dom.createTextNode("\n\n");
          dom.appendChild(el0, el1);
          var el1 = dom.createComment("");
          dom.appendChild(el0, el1);
          var el1 = dom.createTextNode("\n");
          dom.appendChild(el0, el1);
          return el0;
        },
        render: function render(context, env, contextualElement) {
          var dom = env.dom;
          var hooks = env.hooks, get = hooks.get, block = hooks.block;
          dom.detectNamespace(contextualElement);
          var fragment;
          if (env.useFragmentCache && dom.canClone) {
            if (this.cachedFragment === null) {
              fragment = this.build(dom);
              if (this.hasRendered) {
                this.cachedFragment = fragment;
              } else {
                this.hasRendered = true;
              }
            }
            if (this.cachedFragment) {
              fragment = dom.cloneNode(this.cachedFragment, true);
            }
          } else {
            fragment = this.build(dom);
          }
          var element4 = dom.childAt(fragment, [1]);
          var morph0 = dom.createMorphAt(element4,1,1);
          var morph1 = dom.createMorphAt(element4,3,3);
          var morph2 = dom.createMorphAt(fragment,3,3,contextualElement);
          block(env, morph0, context, "if", [get(env, context, "docmeta.isModified")], {}, child0, child1);
          block(env, morph1, context, "if", [get(env, context, "docmeta.isPending")], {}, child2, child3);
          block(env, morph2, context, "if", [get(env, context, "docmeta.isPending")], {}, child4, child5);
          return fragment;
        }
      };
    }());
    return {
      isHTMLBars: true,
      revision: "Ember@1.11.3",
      blockParams: 0,
      cachedFragment: null,
      hasRendered: false,
      build: function build(dom) {
        var el0 = dom.createDocumentFragment();
        var el1 = dom.createTextNode("        ");
        dom.appendChild(el0, el1);
        var el1 = dom.createElement("div");
        dom.setAttribute(el1,"class","col-md-3 col-md-offset-9 sidebar");
        var el2 = dom.createTextNode("\n\n            ");
        dom.appendChild(el1, el2);
        var el2 = dom.createElement("div");
        dom.setAttribute(el2,"class","sline");
        var el3 = dom.createTextNode("\n            ");
        dom.appendChild(el2, el3);
        var el3 = dom.createElement("div");
        dom.setAttribute(el3,"class","icon-createtask");
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode("\n            ");
        dom.appendChild(el2, el3);
        var el3 = dom.createComment("");
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode(" \n            ");
        dom.appendChild(el2, el3);
        dom.appendChild(el1, el2);
        var el2 = dom.createTextNode("\n            ");
        dom.appendChild(el1, el2);
        var el2 = dom.createElement("div");
        dom.setAttribute(el2,"id","list-wrapper");
        var el3 = dom.createTextNode("\n             ");
        dom.appendChild(el2, el3);
        var el3 = dom.createElement("ul");
        var el4 = dom.createTextNode("\n");
        dom.appendChild(el3, el4);
        var el4 = dom.createComment("");
        dom.appendChild(el3, el4);
        var el4 = dom.createTextNode("\n             ");
        dom.appendChild(el3, el4);
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode("\n             ");
        dom.appendChild(el2, el3);
        dom.appendChild(el1, el2);
        var el2 = dom.createTextNode("\n\n        ");
        dom.appendChild(el1, el2);
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n");
        dom.appendChild(el0, el1);
        var el1 = dom.createElement("div");
        dom.setAttribute(el1,"class","site-wrapper");
        dom.setAttribute(el1,"id","main");
        var el2 = dom.createTextNode("\n");
        dom.appendChild(el1, el2);
        var el2 = dom.createComment("");
        dom.appendChild(el1, el2);
        var el2 = dom.createTextNode("\n");
        dom.appendChild(el1, el2);
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n");
        dom.appendChild(el0, el1);
        var el1 = dom.createElement("span");
        dom.setAttribute(el1,"id","nav-icon");
        dom.setAttribute(el1,"class","nav-icon-left");
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n\n");
        dom.appendChild(el0, el1);
        return el0;
      },
      render: function render(context, env, contextualElement) {
        var dom = env.dom;
        var hooks = env.hooks, get = hooks.get, inline = hooks.inline, block = hooks.block, content = hooks.content;
        dom.detectNamespace(contextualElement);
        var fragment;
        if (env.useFragmentCache && dom.canClone) {
          if (this.cachedFragment === null) {
            fragment = this.build(dom);
            if (this.hasRendered) {
              this.cachedFragment = fragment;
            } else {
              this.hasRendered = true;
            }
          }
          if (this.cachedFragment) {
            fragment = dom.cloneNode(this.cachedFragment, true);
          }
        } else {
          fragment = this.build(dom);
        }
        var element5 = dom.childAt(fragment, [1]);
        var morph0 = dom.createMorphAt(dom.childAt(element5, [1]),3,3);
        var morph1 = dom.createMorphAt(dom.childAt(element5, [3, 1]),1,1);
        var morph2 = dom.createMorphAt(dom.childAt(fragment, [3]),1,1);
        inline(env, morph0, context, "input", [], {"type": "text", "id": "new-todo", "class": "form-control", "placeholder": "Enter your title here", "value": get(env, context, "newTitle"), "action": "createDoc"});
        block(env, morph1, context, "each", [get(env, context, "controller")], {"keyword": "docmeta"}, child0, null);
        content(env, morph2, context, "outlet");
        return fragment;
      }
    };
  }()));

});
define('peerdoc-embercli/templates/pd/index', ['exports'], function (exports) {

  'use strict';

  exports['default'] = Ember.HTMLBars.template((function() {
    return {
      isHTMLBars: true,
      revision: "Ember@1.11.3",
      blockParams: 0,
      cachedFragment: null,
      hasRendered: false,
      build: function build(dom) {
        var el0 = dom.createDocumentFragment();
        var el1 = dom.createElement("div");
        dom.setAttribute(el1,"id","welcome");
        var el2 = dom.createTextNode("\n        ");
        dom.appendChild(el1, el2);
        var el2 = dom.createElement("div");
        dom.setAttribute(el2,"class","text-vertical-center");
        var el3 = dom.createTextNode("\n            ");
        dom.appendChild(el2, el3);
        var el3 = dom.createElement("h1");
        var el4 = dom.createTextNode("Start PeerDoc");
        dom.appendChild(el3, el4);
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode("\n            ");
        dom.appendChild(el2, el3);
        var el3 = dom.createElement("h3");
        var el4 = dom.createTextNode("Team20");
        dom.appendChild(el3, el4);
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode("\n            ");
        dom.appendChild(el2, el3);
        var el3 = dom.createElement("br");
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode("\n            ");
        dom.appendChild(el2, el3);
        var el3 = dom.createElement("a");
        dom.setAttribute(el3,"id","findm");
        dom.setAttribute(el3,"href","#about");
        dom.setAttribute(el3,"class","btn btn-dark btn-lg");
        var el4 = dom.createTextNode("Find Out More");
        dom.appendChild(el3, el4);
        dom.appendChild(el2, el3);
        var el3 = dom.createTextNode("\n        ");
        dom.appendChild(el2, el3);
        dom.appendChild(el1, el2);
        var el2 = dom.createTextNode("\n");
        dom.appendChild(el1, el2);
        dom.appendChild(el0, el1);
        var el1 = dom.createTextNode("\n\n");
        dom.appendChild(el0, el1);
        return el0;
      },
      render: function render(context, env, contextualElement) {
        var dom = env.dom;
        dom.detectNamespace(contextualElement);
        var fragment;
        if (env.useFragmentCache && dom.canClone) {
          if (this.cachedFragment === null) {
            fragment = this.build(dom);
            if (this.hasRendered) {
              this.cachedFragment = fragment;
            } else {
              this.hasRendered = true;
            }
          }
          if (this.cachedFragment) {
            fragment = dom.cloneNode(this.cachedFragment, true);
          }
        } else {
          fragment = this.build(dom);
        }
        return fragment;
      }
    };
  }()));

});
define('peerdoc-embercli/tests/adapters/application.jshint', function () {

  'use strict';

  module('JSHint - adapters');
  test('adapters/application.js should pass jshint', function() { 
    ok(true, 'adapters/application.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/app.jshint', function () {

  'use strict';

  module('JSHint - .');
  test('app.js should pass jshint', function() { 
    ok(true, 'app.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/controllers/docmeta.jshint', function () {

  'use strict';

  module('JSHint - controllers');
  test('controllers/docmeta.js should pass jshint', function() { 
    ok(false, 'controllers/docmeta.js should pass jshint.\ncontrollers/docmeta.js: line 7, col 9, \'model\' is defined but never used.\ncontrollers/docmeta.js: line 14, col 9, \'model\' is defined but never used.\n\n2 errors'); 
  });

});
define('peerdoc-embercli/tests/controllers/pd.jshint', function () {

  'use strict';

  module('JSHint - controllers');
  test('controllers/pd.js should pass jshint', function() { 
    ok(false, 'controllers/pd.js should pass jshint.\ncontrollers/pd.js: line 39, col 25, Expected \'!==\' and instead saw \'!=\'.\ncontrollers/pd.js: line 26, col 7, \'$\' is not defined.\n\n2 errors'); 
  });

});
define('peerdoc-embercli/tests/controllers/pd/doc.jshint', function () {

  'use strict';

  module('JSHint - controllers/pd');
  test('controllers/pd/doc.js should pass jshint', function() { 
    ok(true, 'controllers/pd/doc.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/controllers/pd/index.jshint', function () {

  'use strict';

  module('JSHint - controllers/pd');
  test('controllers/pd/index.js should pass jshint', function() { 
    ok(true, 'controllers/pd/index.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/helpers/if-cond.jshint', function () {

  'use strict';

  module('JSHint - helpers');
  test('helpers/if-cond.js should pass jshint', function() { 
    ok(true, 'helpers/if-cond.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/helpers/resolver', ['exports', 'ember/resolver', 'peerdoc-embercli/config/environment'], function (exports, Resolver, config) {

  'use strict';

  var resolver = Resolver['default'].create();

  resolver.namespace = {
    modulePrefix: config['default'].modulePrefix,
    podModulePrefix: config['default'].podModulePrefix
  };

  exports['default'] = resolver;

});
define('peerdoc-embercli/tests/helpers/resolver.jshint', function () {

  'use strict';

  module('JSHint - helpers');
  test('helpers/resolver.js should pass jshint', function() { 
    ok(true, 'helpers/resolver.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/helpers/start-app', ['exports', 'ember', 'peerdoc-embercli/app', 'peerdoc-embercli/router', 'peerdoc-embercli/config/environment'], function (exports, Ember, Application, Router, config) {

  'use strict';



  exports['default'] = startApp;
  function startApp(attrs) {
    var application;

    var attributes = Ember['default'].merge({}, config['default'].APP);
    attributes = Ember['default'].merge(attributes, attrs); // use defaults, but you can override;

    Ember['default'].run(function () {
      application = Application['default'].create(attributes);
      application.setupForTesting();
      application.injectTestHelpers();
    });

    return application;
  }

});
define('peerdoc-embercli/tests/helpers/start-app.jshint', function () {

  'use strict';

  module('JSHint - helpers');
  test('helpers/start-app.js should pass jshint', function() { 
    ok(true, 'helpers/start-app.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/models/doc.jshint', function () {

  'use strict';

  module('JSHint - models');
  test('models/doc.js should pass jshint', function() { 
    ok(true, 'models/doc.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/models/doccg.jshint', function () {

  'use strict';

  module('JSHint - models');
  test('models/doccg.js should pass jshint', function() { 
    ok(true, 'models/doccg.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/models/docdelt.jshint', function () {

  'use strict';

  module('JSHint - models');
  test('models/docdelt.js should pass jshint', function() { 
    ok(true, 'models/docdelt.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/models/docmeta.jshint', function () {

  'use strict';

  module('JSHint - models');
  test('models/docmeta.js should pass jshint', function() { 
    ok(true, 'models/docmeta.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/models/invitation.jshint', function () {

  'use strict';

  module('JSHint - models');
  test('models/invitation.js should pass jshint', function() { 
    ok(true, 'models/invitation.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/models/pd.jshint', function () {

  'use strict';

  module('JSHint - models');
  test('models/pd.js should pass jshint', function() { 
    ok(true, 'models/pd.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/router.jshint', function () {

  'use strict';

  module('JSHint - .');
  test('router.js should pass jshint', function() { 
    ok(true, 'router.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/routes/pd.jshint', function () {

  'use strict';

  module('JSHint - routes');
  test('routes/pd.js should pass jshint', function() { 
    ok(false, 'routes/pd.js should pass jshint.\nroutes/pd.js: line 13, col 43, Expected \'!==\' and instead saw \'!=\'.\nroutes/pd.js: line 102, col 51, Expected \'===\' and instead saw \'==\'.\nroutes/pd.js: line 107, col 71, Missing semicolon.\nroutes/pd.js: line 23, col 9, \'$\' is not defined.\nroutes/pd.js: line 23, col 16, \'$\' is not defined.\nroutes/pd.js: line 24, col 22, \'$\' is not defined.\nroutes/pd.js: line 26, col 17, \'$\' is not defined.\nroutes/pd.js: line 28, col 23, \'$\' is not defined.\nroutes/pd.js: line 32, col 17, \'$\' is not defined.\nroutes/pd.js: line 33, col 23, \'$\' is not defined.\nroutes/pd.js: line 36, col 23, \'$\' is not defined.\nroutes/pd.js: line 39, col 17, \'$\' is not defined.\nroutes/pd.js: line 73, col 12, \'$\' is not defined.\nroutes/pd.js: line 75, col 17, \'$\' is not defined.\nroutes/pd.js: line 80, col 13, \'$\' is not defined.\nroutes/pd.js: line 81, col 17, \'$\' is not defined.\nroutes/pd.js: line 84, col 18, \'$\' is not defined.\nroutes/pd.js: line 88, col 17, \'$\' is not defined.\nroutes/pd.js: line 92, col 14, \'$\' is not defined.\nroutes/pd.js: line 93, col 18, \'$\' is not defined.\nroutes/pd.js: line 94, col 18, \'$\' is not defined.\nroutes/pd.js: line 95, col 18, \'$\' is not defined.\nroutes/pd.js: line 98, col 14, \'$\' is not defined.\nroutes/pd.js: line 100, col 18, \'$\' is not defined.\nroutes/pd.js: line 102, col 21, \'$\' is not defined.\nroutes/pd.js: line 103, col 27, \'$\' is not defined.\nroutes/pd.js: line 107, col 27, \'$\' is not defined.\nroutes/pd.js: line 119, col 30, \'controller\' is defined but never used.\n\n28 errors'); 
  });

});
define('peerdoc-embercli/tests/routes/pd/doc.jshint', function () {

  'use strict';

  module('JSHint - routes/pd');
  test('routes/pd/doc.js should pass jshint', function() { 
    ok(false, 'routes/pd/doc.js should pass jshint.\nroutes/pd/doc.js: line 19, col 65, Expected \'!==\' and instead saw \'!=\'.\nroutes/pd/doc.js: line 64, col 21, \'i\' is already defined.\nroutes/pd/doc.js: line 108, col 34, Expected \'!==\' and instead saw \'!=\'.\nroutes/pd/doc.js: line 27, col 24, \'Quill\' is not defined.\nroutes/pd/doc.js: line 50, col 11, \'$\' is not defined.\nroutes/pd/doc.js: line 149, col 19, \'refresher\' is not defined.\nroutes/pd/doc.js: line 45, col 48, \'source\' is defined but never used.\nroutes/pd/doc.js: line 143, col 28, \'controller\' is defined but never used.\n\n8 errors'); 
  });

});
define('peerdoc-embercli/tests/serializers/docdelt.jshint', function () {

  'use strict';

  module('JSHint - serializers');
  test('serializers/docdelt.js should pass jshint', function() { 
    ok(true, 'serializers/docdelt.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/test-helper', ['peerdoc-embercli/tests/helpers/resolver', 'ember-qunit'], function (resolver, ember_qunit) {

	'use strict';

	ember_qunit.setResolver(resolver['default']);

});
define('peerdoc-embercli/tests/test-helper.jshint', function () {

  'use strict';

  module('JSHint - .');
  test('test-helper.js should pass jshint', function() { 
    ok(true, 'test-helper.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/adapters/application-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('adapter:application', 'ApplicationAdapter', {});

  // Replace this with your real tests.
  ember_qunit.test('it exists', function (assert) {
    var adapter = this.subject();
    assert.ok(adapter);
  });

  // Specify the other units that are required for this test.
  // needs: ['serializer:foo']

});
define('peerdoc-embercli/tests/unit/adapters/application-test.jshint', function () {

  'use strict';

  module('JSHint - unit/adapters');
  test('unit/adapters/application-test.js should pass jshint', function() { 
    ok(true, 'unit/adapters/application-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/controllers/docmeta-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('controller:docmeta', {});

  // Replace this with your real tests.
  ember_qunit.test('it exists', function (assert) {
    var controller = this.subject();
    assert.ok(controller);
  });

  // Specify the other units that are required for this test.
  // needs: ['controller:foo']

});
define('peerdoc-embercli/tests/unit/controllers/docmeta-test.jshint', function () {

  'use strict';

  module('JSHint - unit/controllers');
  test('unit/controllers/docmeta-test.js should pass jshint', function() { 
    ok(true, 'unit/controllers/docmeta-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/controllers/pd-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('controller:pd', {});

  // Replace this with your real tests.
  ember_qunit.test('it exists', function (assert) {
    var controller = this.subject();
    assert.ok(controller);
  });

  // Specify the other units that are required for this test.
  // needs: ['controller:foo']

});
define('peerdoc-embercli/tests/unit/controllers/pd-test.jshint', function () {

  'use strict';

  module('JSHint - unit/controllers');
  test('unit/controllers/pd-test.js should pass jshint', function() { 
    ok(true, 'unit/controllers/pd-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/controllers/pd/doc-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('controller:pd/doc', {});

  // Replace this with your real tests.
  ember_qunit.test('it exists', function (assert) {
    var controller = this.subject();
    assert.ok(controller);
  });

  // Specify the other units that are required for this test.
  // needs: ['controller:foo']

});
define('peerdoc-embercli/tests/unit/controllers/pd/doc-test.jshint', function () {

  'use strict';

  module('JSHint - unit/controllers/pd');
  test('unit/controllers/pd/doc-test.js should pass jshint', function() { 
    ok(true, 'unit/controllers/pd/doc-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/helpers/if-cond-test', ['peerdoc-embercli/helpers/if-cond', 'qunit'], function (if_cond, qunit) {

  'use strict';

  qunit.module('IfCondHelper');

  // Replace this with your real tests.
  qunit.test('it works', function (assert) {
    var result = if_cond.ifCond(42);
    assert.ok(result);
  });

});
define('peerdoc-embercli/tests/unit/helpers/if-cond-test.jshint', function () {

  'use strict';

  module('JSHint - unit/helpers');
  test('unit/helpers/if-cond-test.js should pass jshint', function() { 
    ok(true, 'unit/helpers/if-cond-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/doc-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('doc', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/doc-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/doc-test.js should pass jshint', function() { 
    ok(true, 'unit/models/doc-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/doccg-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('doccg', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/doccg-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/doccg-test.js should pass jshint', function() { 
    ok(true, 'unit/models/doccg-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/docdelt-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('docdelt', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/docdelt-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/docdelt-test.js should pass jshint', function() { 
    ok(true, 'unit/models/docdelt-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/docmeta-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('docmeta', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/docmeta-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/docmeta-test.js should pass jshint', function() { 
    ok(true, 'unit/models/docmeta-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/invitation-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('invitation', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/invitation-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/invitation-test.js should pass jshint', function() { 
    ok(true, 'unit/models/invitation-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/pd-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('pd', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/pd-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/pd-test.js should pass jshint', function() { 
    ok(true, 'unit/models/pd-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/pd/doc-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('pd/doc', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/pd/doc-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models/pd');
  test('unit/models/pd/doc-test.js should pass jshint', function() { 
    ok(true, 'unit/models/pd/doc-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/models/text-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('text', {
    // Specify the other units that are required for this test.
    needs: []
  });

  ember_qunit.test('it exists', function (assert) {
    var model = this.subject();
    // var store = this.store();
    assert.ok(!!model);
  });

});
define('peerdoc-embercli/tests/unit/models/text-test.jshint', function () {

  'use strict';

  module('JSHint - unit/models');
  test('unit/models/text-test.js should pass jshint', function() { 
    ok(true, 'unit/models/text-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/routes/pd-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('route:pd', {});

  ember_qunit.test('it exists', function (assert) {
    var route = this.subject();
    assert.ok(route);
  });

  // Specify the other units that are required for this test.
  // needs: ['controller:foo']

});
define('peerdoc-embercli/tests/unit/routes/pd-test.jshint', function () {

  'use strict';

  module('JSHint - unit/routes');
  test('unit/routes/pd-test.js should pass jshint', function() { 
    ok(true, 'unit/routes/pd-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/routes/pd/doc-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('route:pd/doc', {});

  ember_qunit.test('it exists', function (assert) {
    var route = this.subject();
    assert.ok(route);
  });

  // Specify the other units that are required for this test.
  // needs: ['controller:foo']

});
define('peerdoc-embercli/tests/unit/routes/pd/doc-test.jshint', function () {

  'use strict';

  module('JSHint - unit/routes/pd');
  test('unit/routes/pd/doc-test.js should pass jshint', function() { 
    ok(true, 'unit/routes/pd/doc-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/routes/text-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleFor('route:text', {});

  ember_qunit.test('it exists', function (assert) {
    var route = this.subject();
    assert.ok(route);
  });

  // Specify the other units that are required for this test.
  // needs: ['controller:foo']

});
define('peerdoc-embercli/tests/unit/routes/text-test.jshint', function () {

  'use strict';

  module('JSHint - unit/routes');
  test('unit/routes/text-test.js should pass jshint', function() { 
    ok(true, 'unit/routes/text-test.js should pass jshint.'); 
  });

});
define('peerdoc-embercli/tests/unit/serializers/docdelt-test', ['ember-qunit'], function (ember_qunit) {

  'use strict';

  ember_qunit.moduleForModel('docdelt', {
    // Specify the other units that are required for this test.
    needs: ['serializer:docdelt']
  });

  // Replace this with your real tests.
  ember_qunit.test('it serializes records', function (assert) {
    var record = this.subject();

    var serializedRecord = record.serialize();

    assert.ok(serializedRecord);
  });

});
define('peerdoc-embercli/tests/unit/serializers/docdelt-test.jshint', function () {

  'use strict';

  module('JSHint - unit/serializers');
  test('unit/serializers/docdelt-test.js should pass jshint', function() { 
    ok(true, 'unit/serializers/docdelt-test.js should pass jshint.'); 
  });

});
/* jshint ignore:start */

/* jshint ignore:end */

/* jshint ignore:start */

define('peerdoc-embercli/config/environment', ['ember'], function(Ember) {
  var prefix = 'peerdoc-embercli';
/* jshint ignore:start */

try {
  var metaName = prefix + '/config/environment';
  var rawConfig = Ember['default'].$('meta[name="' + metaName + '"]').attr('content');
  var config = JSON.parse(unescape(rawConfig));

  return { 'default': config };
}
catch(err) {
  throw new Error('Could not read config from meta tag with name "' + metaName + '".');
}

/* jshint ignore:end */

});

if (runningTests) {
  require("peerdoc-embercli/tests/test-helper");
} else {
  require("peerdoc-embercli/app")["default"].create({"name":"peerdoc-embercli","version":"0.0.0.9fcf570f"});
}

/* jshint ignore:end */
//# sourceMappingURL=peerdoc-embercli.map