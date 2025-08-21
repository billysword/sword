<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.10" tiledversion="1.11.2" name="cavern" tilewidth="16" tileheight="16" tilecount="25" columns="25">
 <image source="forest-tiles.png" width="400" height="16"/>
 <tile id="0">
  <properties>
   <property name="kind" value="empty"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
 </tile>
 <tile id="1">
  <properties>
   <property name="kind" value="rock"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="2">
  <properties>
   <property name="kind" value="outer_corner_top_left"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="3">
  <properties>
   <property name="kind" value="right_wall_01"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="4">
  <properties>
   <property name="kind" value="right_wall_02"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="5">
  <properties>
   <property name="kind" value="outer_corner_bottom_left"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="6">
  <properties>
   <property name="kind" value="outer_corner_top_right"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="7">
  <properties>
   <property name="kind" value="left_wall_01"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="8">
  <properties>
   <property name="kind" value="ceiling_01"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="9">
  <properties>
   <property name="kind" value="ceiling_02"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="10">
  <properties>
   <property name="kind" value="barrier_vertical_top"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="11">
  <properties>
   <property name="kind" value="barrier_vertical_bottom"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="12">
  <properties>
   <property name="kind" value="barrier_horizontal_left"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="13">
  <properties>
   <property name="kind" value="barrier_horizontal_right"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="14">
  <properties>
   <property name="kind" value="barrier_cell"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="15">
  <properties>
   <property name="kind" value="barrier_horizontal"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="16">
  <properties>
   <property name="kind" value="barrier_vertical"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="17">
  <properties>
   <property name="kind" value="inner_corner_top_left"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="18">
  <properties>
   <property name="kind" value="inner_corner_top_right"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="19">
  <properties>
   <property name="kind" value="inner_corner_bottom_right"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="20">
  <properties>
   <property name="kind" value="inner_corner_bottom_left"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="21">
  <properties>
   <property name="kind" value="floor_01"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="22">
  <properties>
   <property name="kind" value="floor_02"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="23">
  <properties>
   <property name="kind" value="left_wall_02"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="24">
  <properties>
   <property name="kind" value="outer_corner_bottom_right"/>
   <property name="solid" type="bool" value="true"/>
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
