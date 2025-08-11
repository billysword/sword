<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.10" tiledversion="1.11.2" name="cavern" tilewidth="16" tileheight="16" tilecount="25" columns="25">
 <image source="forest-tiles.png" width="400" height="16"/>
 <tile id="0">
  <properties>
   <property name="kind" value="rock"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="1">
  <properties>
   <property name="kind" value="platform"/>
   <property name="one_way" value="up"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="2">
  <properties>
   <property name="kind" value="slope_up_right"/>
   <property name="slopeB" type="float" value="0"/>
   <property name="slopeM" type="float" value="1"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="3">
  <properties>
   <property name="kind" value="slope_up_left"/>
   <property name="slopeB" type="float" value="16"/>
   <property name="slopeM" type="float" value="-1"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="4">
  <properties>
   <property name="kind" value="spikes"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
 </tile>
 <wangsets>
  <wangset name="CavernTerrain" type="edge" tile="0">
   <wangcolor name="Rock" color="#6b6b6b" tile="0" probability="1"/>
   <wangcolor name="Dirt" color="#8c6239" tile="1" probability="1"/>
   <wangcolor name="Background" color="#2b2b2b" tile="4" probability="1"/>
   <wangtile tileid="0" wangid="1,1,1,1,1,1,1,1"/>
   <wangtile tileid="1" wangid="2,2,2,2,2,2,2,2"/>
   <wangtile tileid="2" wangid="1,0,2,0,1,0,2,0"/>
   <wangtile tileid="3" wangid="2,0,1,0,2,0,1,0"/>
   <wangtile tileid="4" wangid="3,3,3,3,3,3,3,3"/>
  </wangset>
 </wangsets>
</tileset>
