# SPDX-FileCopyrightText: 2026 OpenCHAMI Contributors
# SPDX-License-Identifier: MIT
#
# See `make rpm-build` and docs/RPM_PACKAGING.md for the tag-to-version
# mapping and how the packaged quadlet's image tag is pinned to it.

Name:           inventory-service-quadlet
Version:        %{version}
Release:        %{rel}%{?dist}
Summary:        OpenCHAMI inventory-service Quadlet units

License:        MIT
URL:            https://github.com/OpenCHAMI/inventory-service
Source0:        %{name}-%{version}.tar.gz

BuildArch:      noarch

Requires(post,preun,postun):  systemd

# podman 5.0.0 is the first release whose Quadlet generator enables
# systemd-style *.container.d drop-in directories
Requires:                     podman >= 5.0.0

# NOTE: tokensmith-quadlet is a Suggests rather than a Requires. It is not
# necessary to run the service: auth is optional and tokensmith is simply our
# recommended provider. The OpenCHAMI devs provide separate meta-conf RPMs that
# glue the services together for specific deployment styles, tuning behavior via
# systemd drop-in overrides.
Suggests:                     tokensmith-quadlet >= 0.4.0

%description
Podman Quadlet unit files (container + volume) for running inventory-service
as part of an OpenCHAMI deployment.

%prep
%setup -q

%install
mkdir -p %{buildroot}/usr/share/containers/systemd

grep -q '@IMAGE_TAG@' inventory-service.container
sed "s|@IMAGE_TAG@|v%{version}|" inventory-service.container \
    > %{buildroot}/usr/share/containers/systemd/inventory-service.container
chmod 644 %{buildroot}/usr/share/containers/systemd/inventory-service.container
install -d %{buildroot}/usr/share/containers/systemd/inventory-service.container.d
install -m 644 inventory-service.container.d/10-defaults.conf \
    %{buildroot}/usr/share/containers/systemd/inventory-service.container.d/

install -m 644 inventory-service-data.volume \
    %{buildroot}/usr/share/containers/systemd/inventory-service-data.volume

%files
%license LICENSES/MIT.txt
/usr/share/containers/systemd/inventory-service.container
/usr/share/containers/systemd/inventory-service.container.d
/usr/share/containers/systemd/inventory-service.container.d/10-defaults.conf
/usr/share/containers/systemd/inventory-service-data.volume

%post
# reload systemd so the new Quadlet-generated unit is seen
systemctl daemon-reload || :
if [ $1 -ge 2 ]; then
    systemctl try-restart inventory-service.service || :
fi

%preun
if [ $1 -eq 0 ]; then
    systemctl stop inventory-service.service >/dev/null 2>&1 || :
fi

%postun
# reload systemd so the removed unit is dropped
systemctl daemon-reload || :
