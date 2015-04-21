window.PD = Ember.Application.create({});

//PD.ApplicationAdapter = DS.FixtureAdapter.extend();
PD.DocdeltSerializer = DS.RESTSerializer.extend(DS.EmbeddedRecordsMixin, {
  attrs: {
    doccgs: { embedded: 'always' }
  }
});
/*
DS.RESTSerializer.reopen({
    serializeBelongsTo: function(record, json, relationship) {
        var key = relationship.key,
            belongsTo = Ember.get(record, key);
        key = this.keyForRelationship ? this.keyForRelationship(key, "belongsTo") : key;
        
        if (relationship.options.embedded === 'always') {
            json[key] = belongsTo.serialize();
        }
        else {
            return this._super(record, json, relationship);
        }
    },
    serializeHasMany: function(record, json, relationship) {
        var key = relationship.key,
            hasMany = Ember.get(record, key),
            relationshipType = record.constructor.determineRelationshipType(relationship);
        
        if (relationship.options.embedded === 'always') {
            if (hasMany && relationshipType === 'manyToNone' || relationshipType === 'manyToMany' ||
                relationshipType === 'manyToOne') {
                
                json[key] = [];
                hasMany.forEach(function(item, index){
                  console.log(item);
                    json[key].push(item.record.serialize());
                });
            }
        
        }
        else {
            return this._super(record, json, relationship);
        }
    }
});

*/



PD.ApplicationAdapter = DS.RESTAdapter.extend({
  namespace: 'api',
  host: 'http://localhost:8080',
});

/*PD.ApplicationAdapter.map('PD.Docdelt', {
  doccgs: { embedded: 'always' }
});
*/


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
