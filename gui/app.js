function tremulaInit(){
  $tremulaContainer = $('.tremulaContainer');
  pageCtr = 1

  var tremula = new Tremula();
  var config = tremulaConfigs.default.call(tremula);

  tremula.init($tremulaContainer,config,this);

  var doScrollEvents = function(o){
    if(o.scrollProgress > .7){
      if(!tremula.cache.endOfScrollFlag){
        tremula.cache.endOfScrollFlag = true;
        pageCtr++;
        LoadRoms()
        console.log('Load more')
      }
    }
  }

  tremula.setOnChangePub(doScrollEvents);

  return tremula;
}
function LoadRoms(){
  dataUrl = "http://127.0.0.1:8888/roms/" + pageCtr
  $.getJSON(dataUrl,function(res){

//    var rs = res.roms.filter(function(o,i){  return o.box_art != null });
    tremula.appendData(res.roms,nesAdapter)
  })
}
function applyBoxClick(){
	$('.tremulaContainer').on('tremulaItemSelect',function(gestureEvt,domEvt){
		// console.log(gestureEvt,domEvt)
		var 
		$e = $(domEvt.target);
		t = $e.closest('.gridBox')[0];
		if(t){
			var data = $.data(t).model.model.data;
		}
		if(data){
			console.log(data);
			console.log("Starting", data.ID)
			$.get("http://127.0.0.1:8888/start/"+data.ID)
		}
	})
}
window.loadRoms = LoadRoms

function nesAdapter(data, env){
  this.data = data

console.log(data)
  this.isLastContentBlock   = data.isLastContentBlock||false;
  this.layoutType           = data.layoutType||'tremulaBlockItem';// ['tremulaInline' | 'tremulaBlockItem']
  this.noScaling            = data.noScaling||false;
  this.isFavorite           = data.isFavorite||false;

  this.auxClassList         = data.auxClassList||'';
  this.template = this.data.template||('<img draggable="false" class="moneyShot" onload="imageLoaded(this)" src=""/> <div class="boxLabel nesGame" data-id='+ data.md5 +'">'+ data.sanitized_name +'</div>')

  if( data.box_art && data.box_art[1] != undefined ){
    this.imgUrl = data.box_art[1].src
    this.h = this.height = data.box_art[1].height
    this.w = this.width= data.box_art[1].width
  } else if(data.box_art && data.box_art[0] != undefined){
    this.imgUrl = data.box_art[0].src
    this.h = this.height = data.box_art[0].height
    this.w = this.width = data.box_art[0].width
  }else{
    this.imgUrl = "http://i.imgur.com/nQPpZnW.png"
    this.h = this.height = 900
    this.w = this.width = 654
  }
}


pageCtr = 0

$(document).ready(function(){
  $(".nesGame").live("click", function(e){
	  console.log("clicked", e)
  })
  setTimeout(function(){
    window.tremula = tremulaInit();//does not need to be on the window -- implemented here for convienience.
    LoadRoms()
    applyBoxClick()
  },1000);
});

function tremulaInit(){

  // .tremulaContainer must exist and have actual dimentionality 
  // requires display:block with an explicitly defined H & W
  $tremulaContainer = $('.tremulaContainer');

  //this creates a hook to a new Tremula instance
  var tremula = new Tremula();

  //Create a config object -- this is how most default behaivior is set.
  //see updateConfig(prop_val_object,refreshStreamFlag) method to change properties of a running instance
  var config = {
    //method called after each frame is painted. Passes internal parameter object.
    //see fn definition below
    onChangePub         : doScrollEvents,

    //content/stream data can optionally be passed in on init()
    data                : null,

    // lastContentBlock enables a persistant content block to exist at the end of the stream at all times.
    // Common use case is to target $('.lastContentItem') with a conditional loading spinner when API is churning.
    lastContentBlock    : {
      template :'<div class="lastContentItem"></div>',
      layoutType :'tremulaBlockItem',
      noScaling:true,
      w:300,
      h:300,
      isLastContentBlock:true,
      adapter:tremula.dataAdapters.TremulaItem
    },

    //dafault data adapter method which is called on each data item -- this is used if none is supplied during an import operation
    //enables easy adaptation of arbitrary API data formats -- see flickr example
    adapter             :null,

    //Size of the static axis in pixels
    itemConstraint      :450,

    //Margin in px added to each side of each content item
    itemMargins         :[30,30],//x(l&r),y(t&b) in px

    //Display offset of static axis (static axis is the non-scrolling dimention)
    staticAxisOffset    :0,//px

    //Display offset of scroll axis (this is the amount of scrollable area added before the first content block)
    scrollAxisOffset    :20,//px

    //Sets the scroll axis 'x'|'y'.
    //NOTE: projections generally only work with one scroll axis
    //when changeing this value, make sure to use a compatible projection
    scrollAxis          :'x',//'x'|'y'

    //how many rows (or colums) to display.  note: this is zero based -- so a value of 0 means there will be one row/column
    staticAxisCount     :0,//zero based

    //enables looping with the current seet of results
    isLooping           :false,

    //set this to true to use TremulaJS as a responsive layout machine.
    //when true: ignores user events i.e. touch/pointer/mousewheel events.
    ignoreUserEvents    :false,

    //the grid that will be used to project content
    //NOTE: Generally, this will stay the same and various surface map projections
    //will be used to create various 3d positioning effects
    defaultLayout       :tremula.layouts.basicGridLayout,

    //surfaceMap is the projection/3d-effect which will be used to display grid content
    //following is a list of built-in projections with their corresponding scroll direction
    //NOTE: Using a projection with an incompatible Grid or Grid-Direction will result in-not-so awesome results
    //----------------------
    // (x or y) xyPlain
    // (x) streamHorizontal
    // (y) pinterest
    // (x) mountain
    // (x) turntable
    // (x) enterTheDragon
    // (x) userProjection  [Note: use without namespace like this   surfaceMap:userProjection,  <-- the userProjection is at the bottom of this file...] 
    //----------------------
    surfaceMap          :tremula.projections.streamHorizontal,

    //it does not look like this actually got implemented so, don't worry about it ;)
    itemPreloading      :true,

    //enables the item-level momentum envelope
    itemEasing          :false,

    //if item-level easing is enabled, it will use the following parameters
    //NOTE: this is experimental. This effect can make people queasy.
    itemEasingParams    :{
      touchCurve          :tremula.easings.easeOutCubic,
      swipeCurve          :tremula.easings.easeOutCubic,
      transitionCurve     :tremula.easings.easeOutElastic,
      easeTime            :500,
      springLimit         :40 //in px
    }
  };

  tremula.init($tremulaContainer,config,this);

  return tremula;
}

//This method is called on each paint frame thus enabling low level behaivior control
//it receives a single parameter object of internal instance states
//NOTE: below is a simple example of infinate scrolling where new item
//requests are made when the user scrolls past the existing 70% mark.
//
//Another option here is multiple tremula instancechaining i.e. follow the scroll events of another tremula instance.
//use case of this may be one tremula displaying close up data view while another may be an overview.
function doScrollEvents(o){
  if(o.scrollProgress>.7){
    if(!tremula.cache.endOfScrollFlag){
      tremula.cache.endOfScrollFlag = true;
      pageCtr++;
LoadRoms();
      console.log('END OF SCROLL!')
    }
  }
};
