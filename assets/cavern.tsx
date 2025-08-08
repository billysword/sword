<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.10" tiledversion="1.11.0" name="cavern" tilewidth="16" tileheight="16" tilecount="256" columns="16">
 <image source="cavern.png" width="256" height="256"/>
 <tile id="0">
  <properties>
   <property name="kind" value="air"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
 </tile>
 <tile id="16">
  <properties>
   <property name="kind" value="ground"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="48">
  <properties>
   <property name="kind" value="one_way"/>
   <property name="solid" type="bool" value="true"/>
   <property name="one_way" value="up"/>
  </properties>
 </tile>
 <tile id="64">
  <properties>
   <property name="kind" value="slope"/>
   <property name="solid" type="bool" value="true"/>
   <property name="slopeM" type="float" value="-1"/>
   <property name="slopeB" type="int" value="15"/>
  </properties>
 </tile>
 <tile id="80">
  <properties>
   <property name="kind" value="slope"/>
   <property name="solid" type="bool" value="true"/>
   <property name="slopeM" type="float" value="1"/>
   <property name="slopeB" type="int" value="0"/>
  </properties>
 </tile>
 <tile id="96">
  <properties>
   <property name="kind" value="wall"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="128">
  <properties>
   <property name="kind" value="spikes"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
 </tile>
 <wangsets>
  <wangset name="CavernTerrain" type="corner" tile="16">
   <wangcolor name="Rock" color="#5a4628" tile="16" probability="1"/>
   <wangtile tileid="16" wangid="1,1,1,1,1,1,1,1"/>
   <wangtile tileid="17" wangid="1,1,0,0,1,1,0,0"/>
   <wangtile tileid="18" wangid="0,0,1,1,0,0,1,1"/>
   <wangtile tileid="19" wangid="1,1,1,1,0,0,0,0"/>
   <wangtile tileid="20" wangid="0,0,0,0,1,1,1,1"/>
  </wangset>
 </wangsets>
</tileset>