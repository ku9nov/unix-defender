Name:       unix-defender
Version:    0.0.2
Release:    1
Summary:    Unix-defender

License:    MPL-2.0
URL:        https://example.com
Source0:    ./unix-defender-source.tar.gz
BuildArch:  x86_64

Requires:   at

%global debug_package %{nil}

%description
Unix-defender

%prep
%setup -q


%build


%install
mkdir -p %{buildroot}/etc/unix-defender
mkdir -p %{buildroot}/etc/systemd/system
mkdir -p %{buildroot}/usr/bin
install -m 600 etc/systemd/system/unix-defender.service -t %{buildroot}/etc/systemd/system
install -m 600 etc/unix-defender/.env -t %{buildroot}/etc/unix-defender
install -m 600 etc/unix-defender/iptable-rules.json -t %{buildroot}/etc/unix-defender
cp usr/bin/unix-defender %{buildroot}/usr/bin/unix-defender
install -m 755 usr/bin/unix-defender -t %{buildroot}/usr/bin

%files
%config(noreplace) /etc/unix-defender/.env
%config(noreplace) /etc/unix-defender/iptable-rules.json
/etc/systemd/system/unix-defender.service
/usr/bin/unix-defender


%post
service atd start

%changelog
