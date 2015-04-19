window.PD = Ember.Application.create();

//PD.ApplicationAdapter = DS.FixtureAdapter.extend();
PD.ApplicationAdapter = DS.RESTAdapter.extend({
  namespace: 'api',
  host: 'http://localhost:8080',
});

/*PD.ApplicationAdapter = DS.RESTAdapter.extend({
  namespace: 'api',
  host: 'http://localhost:8080',
  ajax: function(url, method, hash) {
  	hash = hash || {}; // hash may be undefined
    hash.crossDomain = true;
    hash.xhrFields = {withCredentials: true};
    return this._super(url, method, hash);
  }
});*/
