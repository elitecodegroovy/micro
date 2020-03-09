## Build and locally publish dependent libraries

Please build and locally-publish the following dependent libraries.

- pentaho-xul-swt
- org.eclipse.rap.rwt
- org.eclipse.rap.jface
- org.eclipse.rap.fileupload
- org.eclipse.rap.filedialog
- org.eclipse.rap.rwt.testfixture
- pentaho-vfs-browser

### pentaho-xul-swt

```
$ git clone -b webspoon-8.1 https://github.com/HiromuHota/pentaho-commons-xul.git
$ cd pentaho-commons-xul
$ mvn clean install -pl swt
```

### rap

```
$ git clone -b webspoon-3.1-maintenance https://github.com/HiromuHota/rap.git
$ cd rap
$ mvn clean install -N
$ mvn clean install -pl bundles/org.eclipse.rap.rwt -am
$ mvn clean install -pl bundles/org.eclipse.rap.jface -am
$ mvn clean install -pl bundles/org.eclipse.rap.fileupload -am
$ mvn clean install -pl bundles/org.eclipse.rap.filedialog -am
$ mvn clean install -pl tests/org.eclipse.rap.rwt.testfixture -am
```

### pentaho-vfs-browser

```
$ git clone -b webspoon-8.1 https://github.com/HiromuHota/apache-vfs-browser.git
$ cd apache-vfs-browser
$ mvn clean install
```

Remember to build pentaho-vfs-browser after rap due to its dependency on rap.

## Build in the command line

**Make sure patched dependent libraries have been published locally**

Compile `kettle-core`, `kettle-engine`, `kettle-ui-swt`, and `webspoon-security`; and install the packages into the local repository:

```bash
$ git clone -b webspoon-8.1 https://github.com/HiromuHota/pentaho-kettle.git
$ cd pentaho-kettle
$ mvn clean install -pl core,engine,security,ui
```

Build a war file (`spoon.war`):

```bash
$ mvn clean install -pl assemblies/static
$ mvn clean install -pl assemblies/lib
$ mvn clean package -pl assemblies/client
```

## UI testing using Selenium

Currently, only Google Chrome browser has been tested for when running UI test cases.
The tests run in headless mode unless a parameter `-Dheadless.unittest=false` is passed.
To run tests in headless mode, the version of Chrome should be higher than 59.

The default url is `http://localhost:8080/spoon`.
Pass a parameter like below if webSpoon is deployed to a different url.

The following command runs all the unit test cases including UI in non-headless mode.

```
$ mvn clean test -pl integration -Dtest.baseurl=http://localhost:8080/spoon/spoon -Dheadless.unittest=false
```


## 上传到tomcat


```
mkdir ${CATALINA_HOME}/webapps/spoon
unzip -q spoon.war -d ${CATALINA_HOME}/webapps/spoon
```

## 乱码

打开Options，添加参数characterEncoding，设置值为gbk/utf8。

admin/admin123#
guest/guest123#