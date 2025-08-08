<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.10" tiledversion="1.10.2" name="cavern" tilewidth="16" tileheight="16" tilecount="256" columns="16">
  <image source="../../assets/cavern.png" width="256" height="256"/>
  <tile id="0">
    <properties>
      <property name="kind" type="string" value="rock"/>
      <property name="solid" type="bool" value="true"/>
    </properties>
  </tile>
  <tile id="1">
    <properties>
      <property name="kind" type="string" value="platform"/>
      <property name="solid" type="bool" value="true"/>
      <property name="one_way" type="string" value="up"/>
    </properties>
  </tile>
  <tile id="2">
    <properties>
      <property name="kind" type="string" value="slope_up_right"/>
      <property name="solid" type="bool" value="true"/>
      <property name="slopeM" type="float" value="1"/>
      <property name="slopeB" type="float" value="0"/>
    </properties>
  </tile>
  <tile id="3">
    <properties>
      <property name="kind" type="string" value="slope_up_left"/>
      <property name="solid" type="bool" value="true"/>
      <property name="slopeM" type="float" value="-1"/>
      <property name="slopeB" type="float" value="16"/>
    </properties>
  </tile>
  <tile id="4">
    <properties>
      <property name="kind" type="string" value="spikes"/>
      <property name="solid" type="bool" value="false"/>
    </properties>
  </tile>
  <wangsets>
    <wangset name="CavernTerrain" type="edge" tile="0">
      <wangedgecolor name="Rock" color="#6b6b6b" tile="0" probability="1"/>
      <wangedgecolor name="Dirt" color="#8c6239" tile="1" probability="1"/>
      <wangedgecolor name="Background" color="#2b2b2b" tile="4" probability="1"/>
      <wangtile tileid="0" wangid="1,1,1,1,1,1,1,1"/>
      <wangtile tileid="1" wangid="2,2,2,2,2,2,2,2"/>
      <wangtile tileid="2" wangid="1,0,2,0,1,0,2,0"/>
      <wangtile tileid="3" wangid="2,0,1,0,2,0,1,0"/>
      <wangtile tileid="4" wangid="3,3,3,3,3,3,3,3"/>
    </wangset>
  </wangsets>
</tileset>