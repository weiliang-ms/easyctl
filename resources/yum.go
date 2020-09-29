package resources

const CentOSAliBaseYUMContent = "[base]\n" +
	"name=CentOS-$releasever - Base - mirrors.aliyun.com\n" +
	"failovermethod=priority\n" +
	"baseurl=http://mirrors.aliyun.com/centos/$releasever/os/$basearch/\n" +
	"        http://mirrors.aliyuncs.com/centos/$releasever/os/$basearch/\n" +
	"        http://mirrors.cloud.aliyuncs.com/centos/$releasever/os/$basearch/\n" +
	"gpgcheck=1\n" +
	"gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-$releasever\n" +
	"\n" +
	"#released updates\n" +
	"[updates]\n" +
	"name=CentOS-$releasever - Updates - mirrors.aliyun.com\n" +
	"failovermethod=priority\n" +
	"baseurl=http://mirrors.aliyun.com/centos/$releasever/updates/$basearch/\n" +
	"        http://mirrors.aliyuncs.com/centos/$releasever/updates/$basearch/\n" +
	"        http://mirrors.cloud.aliyuncs.com/centos/$releasever/updates/$basearch/\n" +
	"gpgcheck=1\n" +
	"gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-$releasever\n" +
	"\n" +
	"#additional packages that may be useful\n" +
	"[extras]\n" +
	"name=CentOS-$releasever - Extras - mirrors.aliyun.com\n" +
	"failovermethod=priority\n" +
	"baseurl=http://mirrors.aliyun.com/centos/$releasever/extras/$basearch/\n" +
	"        http://mirrors.aliyuncs.com/centos/$releasever/extras/$basearch/\n" +
	"        http://mirrors.cloud.aliyuncs.com/centos/$releasever/extras/$basearch/\n" +
	"gpgcheck=1\ngpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-$releasever\n" +
	"\n" +
	"#additional packages that extend functionality of existing packages\n" +
	"[centosplus]\n" +
	"name=CentOS-$releasever - Plus - mirrors.aliyun.com\n" +
	"failovermethod=priority\n" +
	"baseurl=http://mirrors.aliyun.com/centos/$releasever/centosplus/$basearch/\n" +
	"        http://mirrors.aliyuncs.com/centos/$releasever/centosplus/$basearch/\n" +
	"        http://mirrors.cloud.aliyuncs.com/centos/$releasever/centosplus/$basearch/\n" +
	"gpgcheck=1\nenabled=0\ngpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-$releasever\n" +
	"\n" +
	"#contrib - packages by Centos Users\n" +
	"[contrib]\n" +
	"name=CentOS-$releasever - Contrib - mirrors.aliyun.com\n" +
	"failovermethod=priority\n" +
	"baseurl=http://mirrors.aliyun.com/centos/$releasever/contrib/$basearch/\n" +
	"        http://mirrors.aliyuncs.com/centos/$releasever/contrib/$basearch/\n" +
	"        http://mirrors.cloud.aliyuncs.com/centos/$releasever/contrib/$basearch/\n" +
	"gpgcheck=1\n" +
	"enabled=0\n" +
	"gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-$releaseverF"

const CentOSAliEpelYUMContent = "" +
	"[epel]\n" +
	"name=Extra Packages for Enterprise Linux $releasever - $basearch\n" +
	"baseurl=http://mirrors.aliyun.com/epel/$releasever/$basearch\n" +
	"failovermethod=priority\n" +
	"enabled=1\n" +
	"gpgcheck=0\n" +
	"gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-$releasever\n" +
	"\n" +
	"[epel-debuginfo]\n" +
	"name=Extra Packages for Enterprise Linux $releasever - $basearch - Debug\n" +
	"baseurl=http://mirrors.aliyun.com/epel/$releasever/$basearch/debug\n" +
	"failovermethod=priority\n" +
	"enabled=0\ngpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-$releasever\n" +
	"gpgcheck=0\n" +
	"\n" +
	"[epel-source]\n" +
	"name=Extra Packages for Enterprise Linux $releasever - $basearch - Source\n" +
	"baseurl=http://mirrors.aliyun.com/epel/$releasever/SRPMS\n" +
	"failovermethod=priority\n" +
	"enabled=0\n" +
	"gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-$releasever\n" +
	"gpgcheck=0"

const CentOSLocalYUMContent = "" +
	"[local-repo]\n" +
	"name=local-repo\n" +
	"baseurl=file:///media\n" +
	"enabled=1\n" +
	"gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-$releasever \n" +
	"gpgcheck=0"
