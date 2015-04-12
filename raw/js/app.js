window.PD = Ember.Application.create();

//PD.ApplicationAdapter = DS.FixtureAdapter.extend();
PD.ApplicationAdapter = DS.RESTAdapter.extend({
  namespace: 'api',
  host: 'http://keanelucas.com'
});
