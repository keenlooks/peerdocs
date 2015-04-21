PD.Doc = DS.Model.extend({
  title: DS.attr('string'),
  ctext: DS.attr('string'),
  cursor:DS.attr('number')
});

PD.Docmeta=DS.Model.extend({
  title: DS.attr('string'),
  lastmod: DS.attr('date')
});

PD.Docdelt=DS.Model.extend({
  docid:DS.attr('number'),
  cursor:DS.attr('number'),
  doccgs:DS.hasMany('doccg', {embedded: 'always'}),
});

PD.Doccg=DS.Model.extend({
  location: DS.attr('number'),
  mod: DS.attr('string')
});


PD.Doc.FIXTURES = [
 {
   id: 1,
   title: 'first title',
   ctext: 'this is the content for the first app'
 },
 {
   id: 2,
   title: 'second title',
   ctext: 'this is the content for the second app'
 },
 {
   id: 3,
   title: 'third title',
   ctext: 'this is the content for the third app'
 },
];

PD.Docdelt.FIXTURES=[];

PD.Docmeta.FIXTURES = [
 {
   id: 1,
   title: 'first title',
   lastmod: '2014-05-27T12:54:01'
 },
 {
   id: 2,
   title: 'second title',
   lastmod: '2014-05-27T12:54:01'
 },
 {
   id: 3,
   title: 'third title',
   lastmod: '2014-05-27T12:54:01'
 },
];


/*
{ 
id:1,
delta: [
  {
   id: 1,
   location: number,
   mod: 'text'
 },
 {
   id: 2,
   location: number,
   mod: 'text'
 },
 {
   id: 3,
   location: number,
   mod: 'text'
 },
]

};

*/